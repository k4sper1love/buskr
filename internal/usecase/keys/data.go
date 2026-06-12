package keys

// redis fsm data keys
const (
	// booking
	DataBookingDate = "d_bk_date"
	DataBookingLoc  = "d_bk_loc"
	DataBookingSlot = "d_bk_slt"
	DataBookingMsgID = "d_bk_msg_id"

	// onboarding
	DataOnboardName       = "d_onb_nm"
	DataOnboardNoiseLevel = "d_onb_nlvl"

	// admin location
	DataAdminLocName   = "d_adl_nm"
	DataAdminLocDesc   = "d_adl_dsc"
	DataAdminLocNoise  = "d_adl_no"
	DataAdminLocEditID = "d_adl_eid"
	DataAdminLocLat    = "d_adl_lat"
	DataAdminLocLon    = "d_adl_lon"

	// profile
	DataProfileRequestedNoise = "d_pr_rqn"
)
