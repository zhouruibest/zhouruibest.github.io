// moby/daemon/graphdriver/driver.go
```go
// New creates the driver and initializes it at the specified root.
func New(name string, pg plugingetter.PluginGetter, config Options) (Driver, error) {
	ctx := context.TODO()
	if name != "" { // name = "overlay2"
		log.G(ctx).Infof("[graphdriver] trying configured driver: %s", name)
		if err := checkRemoved(name); err != nil { // 检查是否已经被标记为废弃了
			return nil, err
		}
		return GetDriver(name, pg, config) // 根据name获取驱动
	}

	// Guess for prior driver
	driversMap := scanPriorDrivers(config.Root)
	priorityList := strings.Split(priority, ",")
	log.G(ctx).Debugf("[graphdriver] priority list: %v", priorityList)
	for _, name := range priorityList {
		if _, prior := driversMap[name]; prior {
			// of the state found from prior drivers, check in order of our priority
			// which we would prefer
			driver, err := getBuiltinDriver(name, config.Root, config.DriverOptions, config.IDMap)
			if err != nil {
				// unlike below, we will return error here, because there is prior
				// state, and now it is no longer supported/prereq/compatible, so
				// something changed and needs attention. Otherwise the daemon's
				// images would just "disappear".
				log.G(ctx).Errorf("[graphdriver] prior storage driver %s failed: %s", name, err)
				return nil, err
			}

			// abort starting when there are other prior configured drivers
			// to ensure the user explicitly selects the driver to load
			if len(driversMap) > 1 {
				var driversSlice []string
				for name := range driversMap {
					driversSlice = append(driversSlice, name)
				}

				err = errors.Errorf("%s contains several valid graphdrivers: %s; cleanup or explicitly choose storage driver (-s <DRIVER>)", config.Root, strings.Join(driversSlice, ", "))
				log.G(ctx).Errorf("[graphdriver] %v", err)
				return nil, err
			}

			log.G(ctx).Infof("[graphdriver] using prior storage driver: %s", name)
			return driver, nil
		}
	}

	// If no prior state was found, continue with automatic selection, and pick
	// the first supported, non-deprecated, storage driver (in order of priorityList).
	for _, name := range priorityList {
		driver, err := getBuiltinDriver(name, config.Root, config.DriverOptions, config.IDMap)
		if err != nil {
			if IsDriverNotSupported(err) {
				continue
			}
			return nil, err
		}
		return driver, nil
	}

	// Check all registered drivers if no priority driver is found
	for name, initFunc := range drivers {
		driver, err := initFunc(filepath.Join(config.Root, name), config.DriverOptions, config.IDMap)
		if err != nil {
			if IsDriverNotSupported(err) {
				continue
			}
			return nil, err
		}
		return driver, nil
	}

	return nil, errors.Errorf("no supported storage driver found")
}
```

moby/daemon/graphdriver/driver.go
```go
// GetDriver initializes and returns the registered driver
func GetDriver(name string, pg plugingetter.PluginGetter, config Options) (Driver, error) {
	if initFunc, exists := drivers[name]; exists { // dirvers是一个全局变量，在init()中创建, 同时提供了Register，供具体的驱动调用：将驱动添加到drivers里面
		return initFunc(filepath.Join(config.Root, name), config.DriverOptions, config.IDMap)
	}

	pluginDriver, err := lookupPlugin(name, pg, config)
	if err == nil {
		return pluginDriver, nil
	}
	log.G(context.TODO()).WithError(err).WithField("driver", name).WithField("home-dir", config.Root).Error("Failed to GetDriver graph")
	return nil, ErrNotSupported
}
```


moby/daemon/graphdriver/overlay2/overlay.go
```go
// 注册驱动
func init() {
	graphdriver.Register(driverName, Init)
}


// Init returns the native diff driver for overlay filesystem.
// If overlay filesystem is not supported on the host, the error
// graphdriver.ErrNotSupported is returned.
// If an overlay filesystem is not supported over an existing filesystem then
// the error graphdriver.ErrIncompatibleFS is returned.
func Init(home string, options []string, idMap idtools.IdentityMapping) (graphdriver.Driver, error) {
	opts, err := parseOptions(options)
	if err != nil {
		return nil, err
	}

	// Perform feature detection on /var/lib/docker/overlay2 if it's an existing directory.
	// This covers situations where /var/lib/docker/overlay2 is a mount, and on a different
	// filesystem than /var/lib/docker.
	// If the path does not exist, fall back to using /var/lib/docker for feature detection.
	testdir := home
	if _, err := os.Stat(testdir); os.IsNotExist(err) {
		testdir = filepath.Dir(testdir)
	}

	if err := overlayutils.SupportsOverlay(testdir, true); err != nil {
		logger.Error(err)
		return nil, graphdriver.ErrNotSupported
	}

	fsMagic, err := graphdriver.GetFSMagic(testdir)
	if err != nil {
		return nil, err
	}
	if fsName, ok := graphdriver.FsNames[fsMagic]; ok {
		backingFs = fsName
	}

	supportsDType, err := fs.SupportsDType(testdir)
	if err != nil {
		return nil, err
	}
	if !supportsDType {
		return nil, overlayutils.ErrDTypeNotSupported("overlay2", backingFs)
	}

	usingMetacopy, err := usingMetacopy(testdir)
	if err != nil {
		return nil, err
	}

	cur := idtools.CurrentIdentity()
	dirID := idtools.Identity{
		UID: cur.UID,
		GID: idMap.RootPair().GID,
	}
	if err := idtools.MkdirAllAndChown(home, 0o710, dirID); err != nil {
		return nil, err
	}
	if err := idtools.MkdirAllAndChown(path.Join(home, linkDir), 0o700, cur); err != nil {
		return nil, err
	}

	d := &Driver{
		home:          home,
		idMap:         idMap,
		ctr:           graphdriver.NewRefCounter(graphdriver.NewFsChecker(graphdriver.FsMagicOverlay)),
		supportsDType: supportsDType,
		usingMetacopy: usingMetacopy,
		locker:        locker.New(),
		options:       *opts,
	}

	d.naiveDiff = graphdriver.NewNaiveDiffDriver(d, idMap)

	if backingFs == "xfs" {
		// Try to enable project quota support over xfs.
		if d.quotaCtl, err = quota.NewControl(home); err == nil {
			projectQuotaSupported = true
		} else if opts.quota.Size > 0 {
			return nil, fmt.Errorf("Storage option overlay2.size not supported. Filesystem does not support Project Quota: %v", err)
		}
	} else if opts.quota.Size > 0 {
		// if xfs is not the backing fs then error out if the storage-opt overlay2.size is used.
		return nil, fmt.Errorf("Storage Option overlay2.size only supported for backingFS XFS. Found %v", backingFs)
	}

	// figure out whether "index=off" option is recognized by the kernel
	_, err = os.Stat("/sys/module/overlay/parameters/index")
	switch {
	case err == nil:
		indexOff = "index=off,"
	case os.IsNotExist(err):
		// old kernel, no index -- do nothing
	default:
		logger.Warnf("Unable to detect whether overlay kernel module supports index parameter: %s", err)
	}

	needsUserXattr, err := overlayutils.NeedsUserXAttr(home)
	if err != nil {
		logger.Warnf("Unable to detect whether overlay kernel module needs \"userxattr\" parameter: %s", err)
	}
	if needsUserXattr {
		userxattr = "userxattr,"
	}

	logger.Debugf("backingFs=%s, projectQuotaSupported=%v, usingMetacopy=%v, indexOff=%q, userxattr=%q",
		backingFs, projectQuotaSupported, usingMetacopy, indexOff, userxattr)

	return d, nil
}

```