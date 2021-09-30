package global

const (
	TokenKey = "access-token"
	//CKGToken the key of jwt payload in gin context
	CKGToken = "claims"
)

type PluginStatus int8

const (
	//PluginStatusSubmitted submitted unreviewed
	PluginStatusSubmitted   PluginStatus = 0
	//PluginStatusAudited audited and can publish version
	PluginStatusAudited     PluginStatus = 1
	//PluginStatusPublished the plug-in is available,can start the instance
	PluginStatusPublished   PluginStatus = 2
	//PluginStatusUnavailable the plug-in unavailable,cannot start a new instance
	PluginStatusUnavailable PluginStatus = -1
)

const (
	InsSavedNotRun    InstanceStatus = 0
	InsRunning        InstanceStatus = 1
	InsStartException InstanceStatus = -1
	InsNeedUpdate     InstanceStatus = 2
)

type InstanceStatus int