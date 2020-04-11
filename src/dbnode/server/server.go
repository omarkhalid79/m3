	"os/signal"
	"syscall"
	"github.com/m3db/m3/src/dbnode/storage/namespace"
	xconfig "github.com/m3db/m3/src/x/config"
	"github.com/m3db/m3/src/x/context"
	"github.com/m3db/m3/src/x/ident"
	"github.com/m3db/m3/src/x/instrument"
	xlog "github.com/m3db/m3/src/x/log"
	"github.com/m3db/m3/src/x/pool"
	"github.com/coreos/etcd/embed"
	"github.com/coreos/pkg/capnslog"
	"github.com/jhump/protoreflect/desc"
		logger.Fatalf("could not parse new file mode: %v", err)
		logger.Fatalf("could not parse new directory mode: %v", err)
		logger.Fatalf("could not acquire lock on %s: %v", lockPath, err)
	var schema *desc.MessageDescriptor
	if cfg.Proto != nil {
		logger.Info("Probuf data mode enabled")
		schema, err = proto.ParseProtoSchema(cfg.Proto.SchemaFilePath)
		if err != nil {
			logger.Fatalf("error parsing protobuffer schema: %v", err)
		}
	}

		logger.Fatalf("could not connect to metrics: %v", err)
		logger.Fatalf("could not resolve local host ID: %v", err)
	capnslog.SetGlobalLogLevel(capnslog.WARNING)
		clusters := cfg.EnvironmentConfig.Service.ETCDClusters
				logger.Fatalf("unable to create etcd clusters: %v", err)
			zone := cfg.EnvironmentConfig.Service.Zone
			logger.Infof("using seed nodes etcd cluster: zone=%s, endpoints=%v", zone, endpoints)
			cfg.EnvironmentConfig.Service.ETCDClusters = []etcd.ClusterConfig{etcd.ClusterConfig{
		logger.WithFields(
			xlog.NewField("hostID", hostID),
			xlog.NewField("seedNodeHostIDs", fmt.Sprintf("%v", seedNodeHostIDs)),
		).Info("resolving seed node configuration")
				logger.Fatalf("unable to create etcd config: %v", err)
				logger.Fatalf("could not start embedded etcd: %v", err)
	opts := storage.NewOptions()
	iopts := opts.InstrumentOptions().
		SetLogger(logger).
		SetMetricsScope(scope).
		SetMetricsSamplingRate(cfg.Metrics.SampleRate())
		logger.Warnf("max index query IDs concurrency was not set, falling back to default value")
		logger.Fatalf("unable to start build reporter: %v", err)
		logger.Fatalf("could not construct query cache: %s", err.Error())
		})
		logger.Fatalf("could not set initial runtime options: %v", err)
	mmapCfg := cfg.Filesystem.MmapConfigurationOrDefault()
	shouldUseHugeTLB := mmapCfg.HugeTLB.Enabled
	if shouldUseHugeTLB {
		// Make sure the host supports HugeTLB before proceeding with it to prevent
		// excessive log spam.
		shouldUseHugeTLB, err = hostSupportsHugeTLB()
		if err != nil {
			logger.Fatalf("could not determine if host supports HugeTLB: %v", err)
		}
		if !shouldUseHugeTLB {
			logger.Warnf("host doesn't support HugeTLB, proceeding without it")
		}
	}

		serialize.NewTagDecoderOptions(),
		SetForceBloomFilterMmapMemory(cfg.Filesystem.ForceBloomFilterMmapMemoryOrDefault())
		logger.Fatalf("unknown commit log queue size type: %v",
			cfg.CommitLog.Queue.CalculationType)
			logger.Fatalf("unknown commit log queue channel size type: %v",
				cfg.CommitLog.Queue.CalculationType)
	opts = withEncodingAndPoolingOptions(cfg, logger, schema, opts, cfg.PoolingPolicy)
			SetSegmentReaderPool(opts.SegmentReaderPool()).
			SetIdentifierPool(opts.IdentifierPool())
				retriever := fs.NewBlockRetriever(retrieverOpts, fsopts)
		logger.Fatalf("could not create persist manager: %v", err)
	if cfg.EnvironmentConfig.Static == nil {
			InstrumentOpts: iopts,
			HashingSeed:    cfg.Hashing.Seed,
			logger.Fatalf("could not initialize dynamic config: %v", err)
			logger.Fatalf("could not initialize static config: %v", err)
		runOpts.ClusterClientCh <- envCfg.ClusterClient
	opts = opts.SetNamespaceInitializer(envCfg.NamespaceInitializer)
	topo, err := envCfg.TopologyInitializer.Init()
		logger.Fatalf("could not initialize m3db topology: %v", err)
	m3dbClient, err := cfg.Client.NewAdminClient(
		client.ConfigurationParameters{
			InstrumentOptions: iopts.
				SetMetricsScope(iopts.MetricsScope().SubScope("m3dbclient")),
			TopologyInitializer: envCfg.TopologyInitializer,
		},
		func(opts client.AdminOptions) client.AdminOptions {
			return opts.SetRuntimeOptionsManager(runtimeOptsMgr).(client.AdminOptions)
		},
		func(opts client.AdminOptions) client.AdminOptions {
			return opts.SetContextPool(opts.ContextPool()).(client.AdminOptions)
		},
		func(opts client.AdminOptions) client.AdminOptions {
			return opts.SetOrigin(origin)
		},
		func(opts client.AdminOptions) client.AdminOptions {
			if cfg.Proto != nil {
				return opts.SetEncodingProto(
					schema,
					encoding.NewOptions(),
				).(client.AdminOptions)
			}
			return opts
		},
	)
		logger.Fatalf("could not create m3db client: %v", err)
	// Kick off runtime options manager KV watches
	clientAdminOpts := m3dbClient.Options().(client.AdminOptions)
	kvWatchClientConsistencyLevels(envCfg.KVStore, logger,
		clientAdminOpts, runtimeOptsMgr)
	opts = opts.
		// Feature currently not working.
		SetRepairEnabled(false)
	// Set tchannelthrift options
	ttopts := tchannelthrift.NewOptions().
		SetInstrumentOptions(opts.InstrumentOptions()).
		SetTagEncoderPool(tagEncoderPool).
		SetTagDecoderPool(tagDecoderPool)
	bs, err := cfg.Bootstrap.New(opts, topoMapProvider, origin, m3dbClient)
		logger.Fatalf("could not create bootstrap process: %v", err)
	kvWatchBootstrappers(envCfg.KVStore, logger, timeout, cfg.Bootstrap.Bootstrappers,
				logger.Errorf("updated bootstrapper list is empty")
			updated, err := cfg.Bootstrap.New(opts, topoMapProvider, origin, m3dbClient)
				logger.Errorf("updated bootstrapper list failed: %v", err)
	// Initialize clustered database
	clusterTopoWatch, err := topo.Watch()
		logger.Fatalf("could not create cluster topology watch: %v", err)
	db, err := cluster.NewDatabase(hostID, topo, clusterTopoWatch, opts)
		logger.Fatalf("could not construct database: %v", err)
	if err := db.Open(); err != nil {
		logger.Fatalf("could not open database: %v", err)
	}

	contextPool := opts.ContextPool()

	tchannelOpts := xtchannel.NewDefaultChannelOptions()
	service := ttnode.NewService(db, ttopts)

	tchannelthriftNodeClose, err := ttnode.NewServer(service,
		cfg.ListenAddress, contextPool, tchannelOpts).ListenAndServe()
		logger.Fatalf("could not open tchannelthrift interface on %s: %v",
			cfg.ListenAddress, err)
	defer tchannelthriftNodeClose()
	logger.Infof("node tchannelthrift: listening on %v", cfg.ListenAddress)
	tchannelthriftClusterClose, err := ttcluster.NewServer(m3dbClient,
		cfg.ClusterListenAddress, contextPool, tchannelOpts).ListenAndServe()
		logger.Fatalf("could not open tchannelthrift interface on %s: %v",
			cfg.ClusterListenAddress, err)
	defer tchannelthriftClusterClose()
	logger.Infof("cluster tchannelthrift: listening on %v", cfg.ClusterListenAddress)
	httpjsonNodeClose, err := hjnode.NewServer(service,
		cfg.HTTPNodeListenAddress, contextPool, nil).ListenAndServe()
	if err != nil {
		logger.Fatalf("could not open httpjson interface on %s: %v",
			cfg.HTTPNodeListenAddress, err)
	}
	defer httpjsonNodeClose()
	logger.Infof("node httpjson: listening on %v", cfg.HTTPNodeListenAddress)
	httpjsonClusterClose, err := hjcluster.NewServer(m3dbClient,
		cfg.HTTPClusterListenAddress, contextPool, nil).ListenAndServe()
	if err != nil {
		logger.Fatalf("could not open httpjson interface on %s: %v",
			cfg.HTTPClusterListenAddress, err)
	defer httpjsonClusterClose()
	logger.Infof("cluster httpjson: listening on %v", cfg.HTTPClusterListenAddress)
	if cfg.DebugListenAddress != "" {
		go func() {
			if err := http.ListenAndServe(cfg.DebugListenAddress, nil); err != nil {
				logger.Errorf("debug server could not listen on %s: %v", cfg.DebugListenAddress, err)
			}
		}()
	}
			// Notify on bootstrap chan if specified
		// Bootstrap asynchronously so we can handle interrupt
			logger.Fatalf("could not bootstrap database: %v", err)
		logger.Infof("bootstrapped")
		kvWatchNewSeriesLimitPerShard(envCfg.KVStore, logger, topo,
	// Handle interrupt
	interruptCh := runOpts.InterruptCh
	if interruptCh == nil {
		// Make a noop chan so we can always select
		interruptCh = make(chan error)
	}

	var interruptErr error
	select {
	case err := <-interruptCh:
		interruptErr = err
	case sig := <-interrupt():
		interruptErr = fmt.Errorf("%v", sig)
	}
	logger.Warnf("interrupt: %v", interruptErr)

	// Attempt graceful server close
			logger.Errorf("close database error: %v", err)
	// Wait then close or hard close
		logger.Infof("server closed")
		logger.Errorf("server closed after %s timeout", closeTimeout.String())
func interrupt() <-chan os.Signal {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	return c
}

func bgValidateProcessLimits(logger xlog.Logger) {
		logger.Warnf(`cannot validate process limits: invalid configuration found [%v]`, message)
		logger.WithFields(
			xlog.NewField("url", xdocs.Path("operational_guide/kernel_configuration")),
		).Warnf(`invalid configuration found [%v], refer to linked documentation for more information`, err)
	logger xlog.Logger,
			logger.Warnf("error resolving cluster new series insert limit: %v", err)
		logger.Warnf("unable to set cluster new series insert limit: %v", err)
		logger.Errorf("could not watch cluster new series insert limit: %v", err)
					logger.Warnf("unable to parse new cluster new series insert limit: %v", err)
				logger.Warnf("unable to set cluster new series insert limit: %v", err)
	logger xlog.Logger,
	logger xlog.Logger,
		logger.Errorf("could not resolve KV key %s: %v", key, err)
			logger.Errorf("could not unmarshal KV key %s: %v", key, err)
			logger.Errorf("could not process value of KV key %s: %v", key, err)
			logger.Infof("set KV key %s: %v", key, protoValue.Value)
		logger.Errorf("could not watch KV key %s: %v", key, err)
					logger.Warnf("could not set default for KV key %s: %v", key, err)
				logger.Warnf("could not unmarshal KV key %s: %v", key, err)
				logger.Warnf("could not process change for KV key %s: %v", key, err)
			logger.Infof("set KV key %s: %v", key, protoValue.Value)
	logger xlog.Logger,
		logger.Fatalf("could not watch value for key with KV: %s",
			kvconfig.BootstrapperKey)
				logger.WithFields(
					xlog.NewField("key", kvconfig.BootstrapperKey),
					xlog.NewErrField(err),
				).Error("error converting KV update to string array")
	logger xlog.Logger,
	schema *desc.MessageDescriptor,
		logger.Infof("bytes pool registering bucket capacity=%d, size=%d, "+
			"refillLowWatermark=%f, refillHighWatermark=%f",
			bucket.Capacity, bucket.Size,
			bucket.RefillLowWaterMark, bucket.RefillHighWaterMark)
		logger.Fatalf("unrecognized pooling type: %s", policy.Type)
	logger.Infof("bytes pool %s init", policy.Type)
	bytesPool.Init()
		policy.ContextPool.PoolPolicy,
	contextPool := context.NewPool(context.NewOptions().
		SetFinalizerPoolOptions(closersPoolOpts).
		SetMaxPooledFinalizerCapacity(policy.ContextPool.MaxFinalizerCapacityOrDefault()))
		SetSegmentReaderPool(segmentReaderPool)
		if schema != nil {
			enc.SetSchema(schema)
	iteratorPool.Init(func(r io.Reader) encoding.ReaderIterator {
		if schema != nil {
			return proto.NewIterator(r, schema, encodingOpts)
	multiIteratorPool.Init(func(r io.Reader) encoding.ReaderIterator {
		iter.Reset(r)
		SetWriteBatchPool(writeBatchPool)
		return block.NewDatabaseBlock(time.Time{}, 0, ts.Segment{}, blockOpts)
				SetInstrumentOptions(iopts)).
		SetAggregateResultsPool(aggregateQueryResultsPool)
	defer mmap.Munmap(withHugeTLB.Result)
	defer mmap.Munmap(withoutHugeTLB.Result)