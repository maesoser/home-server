package virgin

var DownChannelEntries = []string{
	"down_chann_ID",
	"down_chann_Freq",
	"down_chann_Width",
	"down_chann_Modulation",
	"down_chann_Interleave",
	"down_chann_Power",
	"down_chann_Annex",
	"down_chann_StorageType",
}

var WiFiClientEntries = []string{
	"wifi_client_Index",
	"wifi_client_IPAddrType",
	"wifi_client_IPAddr",
	"wifi_client_IPAddrTextual",
	"wifi_client_HostName",
	"wifi_client_MAC",
	"wifi_client_MACMfg",
	"wifi_client_Status",
	"wifi_client_FirstSeen",
	"wifi_client_LastSeen",
	"wifi_client_IdleTime",
	"wifi_client_InNetworkTime",
	"wifi_client_State",
	"wifi_client_Flags",
	"wifi_client_TxPkts",
	"wifi_client_TxFailures",
	"wifi_client_RxUnicastPkts",
	"wifi_client_RxMulticastPkts",
	"wifi_client_LastTxPktRate",
	"wifi_client_LastRxPktRate",
	"wifi_client_RateSet",
	"wifi_client_RSSI",
}

var UpChannelEntries = []string{
	"up_channel_ID",
	"up_channel_Freq",
	"up_channel_Width",
	"up_channel_ModulationProfile",
	"up_channel_SlotSize",
	"up_channel_TxTimingOffset",
	"up_channel_RangingBackoffStart",
	"up_channel_RangingBackoffEnd",
	"up_channel_TxBackoffStart",
	"up_channel_TxBackoffEnd",
	"up_channel_ScdmaActiveCodes",
	"up_channel_ScdmaCodesPerSlot",
	"up_channel_ScdmaFrameSize",
	"up_channel_ScdmaHoppingSeed",
	"up_channel_Type",
	"up_channel_CloneFrom",
	"up_channel_Update",
	"up_channel_Status",
	"up_channel_PreEqEnable",
}

var SignalQualityEntry = []string{
	"siqnal_quality_IncludesContention",
	"siqnal_quality_Unerroreds",
	"siqnal_quality_Correcteds",
	"siqnal_quality_Uncorrectables",
	"siqnal_quality_SignalNoise",
	"siqnal_quality_Microreflections",
	"siqnal_quality_EqualizationData",
	"siqnal_quality_ExtUnerroreds",
	"siqnal_quality_ExtCorrecteds",
	"siqnal_quality_ExtUncorrectables",
}

var DevEventEntry = []string{
	"docsDevEvIndex",
	"docsDevEvFirstTime",
	"docsDevEvLastTime",
	"docsDevEvCounts",
	"docsDevEvLevel",
	"docsDevEvId",
	"docsDevEvText",
}

var FwEventEntry = []string{
	"FwEvIndex",
	"FwEvTime",
	"FwEvText",
}
