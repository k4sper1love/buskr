package keys

// common
const (
	TextCommonErrGeneral      = "common.err.general"
	TextCommonBtnBack         = "common.btn.back"
	TextCommonBtnMenu         = "common.btn.menu"
	TextCommonBtnCancel       = "common.btn.cancel"
	TextCommonBtnOpenMap      = "common.btn.open_map"
	TextCommonBtnLastLoc      = "common.btn.last_loc"
	TextCommonBtnSchedule     = "common.btn.schedule"
	TextCommonLblNoise        = "common.lbl.noise"
	TextCommonLblNoiseLight   = "common.lbl.noise_light"
	TextCommonLblNoiseMedium  = "common.lbl.noise_medium"
	TextCommonLblNoiseHard    = "common.lbl.noise_hard"
	TextCommonLblNoiseNone    = "common.lbl.noise_none"
	TextCommonLblRoleMusician = "common.lbl.role_musician"
	TextCommonLblRoleAdmin    = "common.lbl.role_admin"

	TextCommonWeekdaySun  = "common.weekday.sun"
	TextCommonWeekdayMon  = "common.weekday.mon"
	TextCommongWeekdayTue = "common.weekday.tue"
	TextCommonWeekdayWed  = "common.weekday.wed"
	TextCommonWeekdayThu  = "common.weekday.thu"
	TextCommonWeekdayFri  = "common.weekday.fri"
	TextCommonWeekdaySat  = "common.weekday.sat"

	TextCommonWeekdaySunFull = "common.weekday_full.sun"
	TextCommonWeekdayMonFull = "common.weekday_full.mon"
	TextCommonWeekdayTueFull = "common.weekday_full.tue"
	TextCommonWeekdayWedFull = "common.weekday_full.wed"
	TextCommonWeekdayThuFull = "common.weekday_full.thu"
	TextCommonWeekdayFriFull = "common.weekday_full.fri"
	TextCommonWeekdaySatFull = "common.weekday_full.sat"

	TextCommonMonthJan = "common.month.jan"
	TextCommonMonthFeb = "common.month.feb"
	TextCommonMonthMar = "common.month.mar"
	TextCommonMonthApr = "common.month.apr"
	TextCommonMonthMay = "common.month.may"
	TextCommonMonthJun = "common.month.jun"
	TextCommonMonthJul = "common.month.jul"
	TextCommonMonthAug = "common.month.aug"
	TextCommonMonthSep = "common.month.sep"
	TextCommonMonthOct = "common.month.oct"
	TextCommonMonthNov = "common.month.nov"
	TextCommonMonthDec = "common.month.dec"
)

// auth
const (
	TextAuthGuestTitle        = "auth.guest.title"
	TextAuthGuestBtnApply     = "auth.guest.btn_apply"
	TextAuthPendingTitle      = "auth.pending.title"
	TextAuthBannedTitle       = "auth.banned.title"
	TextAuthActiveTitle       = "auth.active.title"
	TextAuthActiveBtnBook     = "auth.active.btn_book"
	TextAuthActiveBtnBookings = "auth.active.btn_bookings"
	TextAuthActiveBtnProfile  = "auth.active.btn_profile"
	TextAuthActiveBtnAdmin    = "auth.active.btn_admin"
	TextAuthInvitedPromptName     = "auth.invited.prompt_name"
	TextAuthInvitedErrNotFound    = "auth.invited.err_not_found"
	TextAuthInvitedErrAlreadyUsed = "auth.invited.err_already_used"
	TextAuthInvitedErrExpired     = "auth.invited.err_expired"
	TextAuthInvitedErrActive      = "auth.invited.err_active"
	TextAuthInvitedErrGeneric     = "auth.invited.err_generic"
)

// group
const (
	TextGroupUnathorizedChat = "group.unauthorized_chat"
	TextGroupWelcomeUser     = "group.welcome_user"
)

// onboarding
const (
	TextOnboardStep1PromptName  = "onboard.step1.prompt_name"
	TextOnboardStep2PromptNoise = "onboard.step2.prompt_noise"
	TextOnboardStep3PromptMedia = "onboard.step3.prompt_media"
	TextOnboardBtnSolo          = "onboard.btn_solo"
	TextOnboardBtnGroup         = "onboard.btn_group"
	TextOnboardMsgSuccess       = "onboard.msg_success"
	TextOnboardMsgErr           = "onboard.msg_err"
	TextOnboardMsgCancel        = "onboard.msg_cancel"
	TextOnboardBtnSkip          = "onboard.btn_skip"
)

// admin
const (
	// Panel
	TextAdminPanelTitle      = "admin.panel.title"
	TextAdminPanelBtnInvites = "admin.panel.btn_invites"
	TextAdminPanelBtnLocs    = "admin.panel.btn_locs"
	TextAdminPanelBtnUsers   = "admin.panel.btn_users"

	// Invites
	TextAdminInvTitle   = "admin.inv.title"
	TextAdminInvCreated = "admin.inv.created"

	// Moderation
	TextAdminModBtnApprove    = "admin.mod.btn_approve"
	TextAdminModBtnReject     = "admin.mod.btn_reject"
	TextAdminModAppTitle      = "admin.mod.app_title"
	TextAdminModAppLink       = "admin.mod.app_link"
	TextAdminModNoiseTitle    = "admin.mod.noise_title"
	TextAdminModMsgApprSfx    = "admin.mod.msg_appr_sfx"
	TextAdminModMsgApprCb     = "admin.mod.msg_appr_cb"
	TextAdminModMsgApprNotify = "admin.mod.msg_appr_notify"
	TextAdminModMsgRejSfx     = "admin.mod.msg_rej_sfx"
	TextAdminModMsgRejCb      = "admin.mod.msg_rej_cb"
	TextAdminModMsgRejNotify  = "admin.mod.msg_rej_notify"
	TextAdminModMsgUpgSfx     = "admin.mod.msg_upg_sfx"
	TextAdminModMsgUpgCb      = "admin.mod.msg_upg_cb"
	TextAdminModMsgUpgNotify  = "admin.mod.msg_upg_notify"
	TextAdminModAppNoMedia    = "admin.mod.app_no_media"

	// Locations
	TextAdminLocsTitle           = "admin.locs.title"
	TextAdminLocsBtnAdd          = "admin.locs.btn_add"
	TextAdminLocsAddStep1        = "admin.locs.add_step1"
	TextAdminLocsAddStep2        = "admin.locs.add_step2"
	TextAdminLocsAddStep3        = "admin.locs.add_step3"
	TextAdminLocsAddStep4        = "admin.locs.add_step4"
	TextAdminLocsMsgAddSuccess   = "admin.locs.msg_add_success"
	TextAdminLocsMsgAddErr       = "admin.locs.msg_add_err"
	TextAdminLocsMsgCancel       = "admin.locs.msg_cancel"
	TextAdminLocsMsgNotFound     = "admin.locs.msg_not_found"
	TextAdminLocsBtnList         = "admin.locs.btn_list"
	TextAdminLocsBtnDelete       = "admin.locs.btn_delete"
	TextAdminLocsDelConfirmTitle = "admin.locs.del_confirm_title"
	TextAdminLocsDelConfirmBtn   = "admin.locs.del_confirm_btn"
	TextAdminLocsDetails         = "admin.locs.details"
	TextAdminLocsBtnEnable       = "admin.locs.btn_enable"
	TextAdminLocsBtnDisable      = "admin.locs.btn_disable"

	TextAdminLocsBtnEdit        = "admin.locs.btn_edit"
	TextAdminLocsEditTitle      = "admin.locs.edit_title"
	TextAdminLocsEditBtnName    = "admin.locs.edit_btn_name"
	TextAdminLocsEditBtnDesc    = "admin.locs.edit_btn_desc"
	TextAdminLocsEditBtnGeo     = "admin.locs.edit_btn_geo"
	TextAdminLocsEditStepName   = "admin.locs.edit_step_name"
	TextAdminLocsEditStepDesc   = "admin.locs.edit_step_desc"
	TextAdminLocsEditStepGeo    = "admin.locs.edit_step_geo"
	TextAdminLocsEditMsgSuccess = "admin.locs.edit_msg_success"
	TextAdminLocsEditMsgErr     = "admin.locs.edit_msg_err"
	TextAdminLocsEditBtnNoise   = "admin.locs.edit.btn_noise"
	TextAdminLocsEditNoiseTitle = "admin.locs.edit.noise_title"

	TextAdminLocsScheduleEmpty             = "admin.locs.schedule_empty"
	TextAdminLocsScheduleLblTodayActive    = "admin.locs.schedule.lbl_today_active"
	TextAdminLocsScheduleLblTomorrowActive = "admin.locs.schedule.lbl_tomorrow_active"
	TextAdminLocsScheduleLblOtherActive    = "admin.locs.schedule.lbl_other_active"
	TextAdminLocsScheduleEmptyForDay       = "admin.locs.schedule_empty_for_day"
	TextAdminLocsNearbyWarn                = "admin.locs.nearby_warning"
	TextAdminLocsNearbyConfirm             = "admin.locs.nearby_confirm"
	TextAdminLocsMapTitle                  = "admin.locs.map_title"
	TextAdminLocsMapEmpty                  = "admin.locs.map_empty"
	TextAdminLocsErrDeleteHasBookings      = "admin.locs.err.delete_has_bookings"

	// Users
	TextAdminUsersListTitle          = "admin.users.list_title"
	TextAdminUsersBtnSearch          = "admin.users.btn_search"
	TextAdminUsersSearchQueryPrompt  = "admin.users.search_query_prompt"
	TextAdminUsersSearchResultsTitle = "admin.users.search_results_title"
	TextAdminUsersPromptSearch       = "admin.users.prompt_search"
	TextAdminUsersMsgInvalid         = "admin.users.msg_invalid"
	TextAdminUsersMsgNotFound        = "admin.users.msg_not_found"
	TextAdminUsersSearchResult       = "admin.users.search_result"
	TextAdminUsersBtnChangeNoise     = "admin.users.btn_change_noise"
	TextAdminUsersBtnSortDate        = "admin.users.btn_sort_date"
	TextAdminUsersBtnSortKarmaAsc    = "admin.users.btn_sort_karma_asc"
	TextAdminUsersBtnSortRole        = "admin.users.btn_sort_role"
	TextAdminUsersBtnSortName        = "admin.users.btn_sort_name"
	TextAdminUsersNoiseMenuTitle     = "admin.users.noise_menu_title"
	TextAdminUsersBtnBan             = "admin.users.btn_ban"
	TextAdminUsersBtnUnban           = "admin.users.btn_unban"
	TextAdminUsersBtnPromote         = "admin.users.btn_promote"
	TextAdminUsersBtnDemote          = "admin.users.btn_demote"
	TextAdminUsersMsgIsAdmin         = "admin.users.msg_is_admin"
	TextAdminUsersMsgBanSuccess      = "admin.users.msg_ban_success"
	TextAdminUsersMsgUnbanSuccess    = "admin.users.msg_unban_success"
)

// profile
const (
	TextProfileMainTitle          = "profile.main.title"
	TextProfileMainLblKarmaGood   = "profile.main.lbl_karma_good"
	TextProfileMainLblKarmaBad    = "profile.main.lbl_karma_bad"
	TextProfileMainBtnEditName    = "profile.main.btn_edit_name"
	TextProfileMainBtnUpgNoise    = "profile.main.btn_upg_noise"
	TextProfileEditNamePrompt     = "profile.edit_name.prompt"
	TextProfileEditNameMsgSuccess = "profile.edit_name.msg_success"
	TextProfileEditNameMsgErr     = "profile.edit_name.msg_err"
	TextProfileUpgMsgSuccess      = "profile.upg.msg_success"
	TextProfileUpgAlreadyPending  = "profile.upg.msg_already_pending"
	TextProfileUpgPromptReason    = "profile.upg.prompt_reason"
	TextProfileUpgBtnSkip         = "profile.upg.btn_skip"
)

// booking
const (
	// create flow
	TextBookCreatePromptDate   = "book.create.prompt_date"
	TextBookCreateLblToday     = "book.create.lbl_today"
	TextBookCreateLblTomorrow  = "book.create.lbl_tomorrow"
	TextBookCreateLblOther     = "book.create.lbl_other"
	TextBookCreateMsgNoLocs    = "book.create.msg_no_locs"
	TextBookCreateBtnDates     = "book.create.btn_dates"
	TextBookCreatePromptLoc    = "book.create.prompt_loc"
	TextBookCreatePromptLocErr = "book.create.prompt_loc_err"
	TextBookCreateMsgNoSlots   = "book.create.msg_no_slots"
	TextBookCreateBtnLocs      = "book.create.btn_locs"
	TextBookCreatePromptSlot   = "book.create.prompt_slot"
	TextBookCreateBtnSlots     = "book.create.btn_slots"
	TextBookCreatePromptDur    = "book.create.prompt_dur"
	TextBookCreateLblHours     = "book.create.lbl_hours"
	TextBookCreateMsgSuccess   = "book.create.msg_success"
	TextBookCreateMsgSuccessCb = "book.create.msg_success_cb"
	TextBookCreateMsgCancel    = "book.create.msg_cancel"
	TextBookCreateBtn2GIS      = "book.create.btn_2gis"
	TextBookCreateBtnYandex    = "book.create.btn_yandex"

	// list flow
	TextBookListTitle    = "book.list.title"
	TextBookListMsgEmpty = "book.list.msg_empty"

	// details flow
	TextBookDetTitle            = "book.det.title"
	TextBookDetLblPending       = "book.det.lbl_pending"
	TextBookDetLblActive        = "book.det.lbl_active"
	TextBookDetLblCompleted     = "book.det.lbl_completed"
	TextBookDetLblCancelled     = "book.det.lbl_cancelled"
	TextBookDetLblNoshow        = "book.det.lbl_noshow"
	TextBookDetLblUnknown       = "book.det.lbl_unknown"
	TextBookDetBtnCheckin       = "book.det.btn_checkin"
	TextBookDetBtnCancel        = "book.det.btn_cancel"
	TextBookDetBtnList          = "book.det.btn_list"
	TextBookDetMsgCancelSuccess = "book.det.msg_cancel_success"
	TextBookDetMsgCheckinSfx    = "book.det.msg_checkin_sfx"
	TextBookDetMsgCheckinCb     = "book.det.msg_checkin_cb"
	TextBookDetMsgGrabTitle     = "book.det.msg_grab_title"

	// schedule
	TextBookSchedulePromptLoc   = "book.schedule.prompt_loc"
	TextBookScheduleEmptyForDay = "book.schedule.empty_for_day"

	// errors
	TextBookErrSlotTaken        = "book.err.slot_taken"
	TextBookErrNoisyNeighbor    = "book.err.noisy_neighbor"
	TextBookErrNoiseExceeded    = "book.err.noise_exceeded"
	TextBookErrMaxActive        = "book.err.max_active"
	TextBookErrMaxPerLoc        = "book.err.max_per_loc"
	TextBookErrTooFarInFuture   = "book.err.too_far"
	TextBookErrTimeOverlap      = "book.err.time_overlap"
	TextBookErrInvalidTime      = "book.err.invalid_time"
	TextBookErrInvalidStatus    = "book.err.invalid_status"
	TextBookErrHotSlotsDisabled = "book.err.hot_slots_disabled"
)

// worker
const (
	TextWorkerReminderMsg        = "worker.reminder.msg"
	TextWorkerReminderBtn        = "worker.reminder.btn"
	TextWorkerCheckinFailMsg     = "worker.checkin.fail_msg"
	TextWorkerCheckinFailSoftMsg = "worker.checkin.fail_soft_msg"
	TextWorkerHotSpotMsg         = "worker.hotspot.msg"
	TextWorkerHotSpotBtn         = "worker.hotspot.btn"
	TextWorkerHotSpotDismissBtn  = "worker.hotspot.dismiss_btn"
)
