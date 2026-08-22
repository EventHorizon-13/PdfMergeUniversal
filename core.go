//go:build windows

package main

import (
	"sync"
	"syscall"
)

const (
	appTitle = "PDF 合并工具 v3.1.4"

	WS_OVERLAPPED       = 0x00000000
	WS_CAPTION          = 0x00C00000
	WS_SYSMENU          = 0x00080000
	WS_THICKFRAME       = 0x00040000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_OVERLAPPEDWINDOW = WS_OVERLAPPED | WS_CAPTION | WS_SYSMENU | WS_THICKFRAME | WS_MINIMIZEBOX | WS_MAXIMIZEBOX
	WS_CHILD            = 0x40000000
	WS_VISIBLE          = 0x10000000
	WS_TABSTOP          = 0x00010000
	WS_BORDER           = 0x00800000
	WS_VSCROLL          = 0x00200000
	WS_HSCROLL          = 0x00100000
	WS_EX_CLIENTEDGE    = 0x00000200

	SS_WHITERECT   = 0x00000006
	SS_CENTER      = 0x00000001
	SS_CENTERIMAGE = 0x00000200
	SS_OWNERDRAW   = 0x0000000D

	ES_AUTOHSCROLL   = 0x0080
	BS_PUSHBUTTON    = 0x00000000
	BS_DEFPUSHBUTTON = 0x00000001
	BS_AUTOCHECKBOX  = 0x00000003
	BS_OWNERDRAW     = 0x0000000B

	LVS_REPORT           = 0x0001
	LVS_SHOWSELALWAYS    = 0x0008
	LVS_EX_GRIDLINES     = 0x00000001
	LVS_EX_FULLROWSELECT = 0x00000020
	LVS_EX_DOUBLEBUFFER  = 0x00010000

	LVCF_FMT      = 0x0001
	LVCF_WIDTH    = 0x0002
	LVCF_TEXT     = 0x0004
	LVCFMT_LEFT   = 0x0000
	LVCFMT_RIGHT  = 0x0001
	LVIF_TEXT     = 0x0001
	LVIS_FOCUSED  = 0x0001
	LVIS_SELECTED = 0x0002
	LVNI_SELECTED = 0x0002

	LVM_FIRST                    = 0x1000
	LVM_GETITEMCOUNT             = LVM_FIRST + 4
	LVM_DELETEALLITEMS           = LVM_FIRST + 9
	LVM_GETNEXTITEM              = LVM_FIRST + 12
	LVM_ENSUREVISIBLE            = LVM_FIRST + 19
	LVM_SETITEMSTATE             = LVM_FIRST + 43
	LVM_SETEXTENDEDLISTVIEWSTYLE = LVM_FIRST + 54
	LVM_INSERTITEMW              = LVM_FIRST + 77
	LVM_SETITEMW                 = LVM_FIRST + 76
	LVM_INSERTCOLUMNW            = LVM_FIRST + 97

	WM_CREATE         = 0x0001
	WM_DESTROY        = 0x0002
	WM_SIZE           = 0x0005
	WM_COMMAND        = 0x0111
	WM_NOTIFY         = 0x004E
	WM_DRAWITEM       = 0x002B
	WM_CTLCOLORSTATIC = 0x0138
	WM_GETMINMAXINFO  = 0x0024
	WM_DROPFILES      = 0x0233
	WM_TIMER          = 0x0113
	WM_SETFONT        = 0x0030
	WM_USER           = 0x0400
	WM_APP            = 0x8000

	EN_CHANGE  = 0x0300
	BN_CLICKED = 0

	SW_SHOW       = 5
	SW_SHOWNORMAL = 1
	SW_RESTORE    = 9
	SWP_NOMOVE    = 0x0002
	SWP_NOSIZE    = 0x0001

	CW_USEDEFAULT = ^uintptr(0x7fffffff)

	IDC_ARROW       = 32512
	IDI_APPLICATION = 32512
	IMAGE_ICON      = 1
	LR_LOADFROMFILE = 0x0010
	LR_DEFAULTSIZE  = 0x0040

	COLOR_WINDOW  = 5
	COLOR_BTNFACE = 15

	ODS_SELECTED  = 0x0001
	ODS_DISABLED  = 0x0004
	ODS_FOCUS     = 0x0010
	DT_CENTER     = 0x00000001
	DT_VCENTER    = 0x00000004
	DT_SINGLELINE = 0x00000020
	TRANSPARENT   = 1
	PS_SOLID      = 0

	MB_OK              = 0x00000000
	MB_ICONERROR       = 0x00000010
	MB_ICONWARNING     = 0x00000030
	MB_ICONINFORMATION = 0x00000040
	MB_YESNO           = 0x00000004
	IDYES              = 6

	PBM_SETRANGE32 = WM_USER + 6
	PBM_SETPOS     = WM_USER + 2
	PBS_SMOOTH     = 0x01

	NM_DBLCLK = -3
	NM_RCLICK = -5

	MF_STRING       = 0x00000000
	MF_SEPARATOR    = 0x00000800
	MF_POPUP        = 0x00000010
	TPM_RIGHTBUTTON = 0x0002
	TPM_RETURNCMD   = 0x0100

	OFN_EXPLORER         = 0x00080000
	OFN_FILEMUSTEXIST    = 0x00001000
	OFN_PATHMUSTEXIST    = 0x00000800
	OFN_ALLOWMULTISELECT = 0x00000200
	OFN_OVERWRITEPROMPT  = 0x00000002

	BIF_RETURNONLYFSDIRS = 0x00000001
	BIF_NEWDIALOGSTYLE   = 0x00000040

	FVIRTKEY  = 0x01
	FCONTROL  = 0x08
	FALT      = 0x10
	VK_RETURN = 0x0D
	VK_DELETE = 0x2E
	VK_UP     = 0x26
	VK_DOWN   = 0x28

	ERROR_ALREADY_EXISTS    = 183
	SPI_GETICONTITLELOGFONT = 0x001F

	ID_LIST       = 1001
	ID_ADD        = 1002
	ID_DELETE     = 1003
	ID_UP         = 1004
	ID_DOWN       = 1005
	ID_NATSORT    = 1006
	ID_REVERSE    = 1007
	ID_EXPORT_CSV = 1008
	ID_IMPORT_CSV = 1009
	ID_CLEAR      = 1010
	ID_BROWSE_DIR = 1011
	ID_NAME       = 1012
	ID_LOCATE     = 1013
	ID_MERGE      = 1014
	ID_SELECT_ALL = 1015
	ID_TOOLS      = 1016

	ID_MENU_ADD       = 2001
	ID_MENU_CLEAR     = 2002
	ID_MENU_EXIT      = 2003
	ID_MENU_INSTALL   = 2010
	ID_MENU_UNINSTALL = 2011
	ID_MENU_QPDF      = 2012

	ID_CTX_OPEN   = 2101
	ID_CTX_FOLDER = 2102
	ID_CTX_TOP    = 2103
	ID_CTX_BOTTOM = 2104
	ID_CTX_UP     = 2105
	ID_CTX_DOWN   = 2106
	ID_CTX_DELETE = 2107

	TIMER_QUEUE = 1

	MSG_REFRESH    = WM_APP + 101
	MSG_MERGE_DONE = WM_APP + 102
)

type POINT struct{ X, Y int32 }
type RECT struct{ Left, Top, Right, Bottom int32 }
type DRAWITEMSTRUCT struct {
	CtlType uint32
	CtlID uint32
	ItemID uint32
	ItemAction uint32
	ItemState uint32
	HwndItem uintptr
	HDC uintptr
	RcItem RECT
	ItemData uintptr
}
type MSG struct {
	Hwnd uintptr
	Message uint32
	WParam, LParam uintptr
	Time uint32
	Pt POINT
}
type WNDCLASSEX struct {
	CbSize uint32
	Style uint32
	LpfnWndProc uintptr
	CbClsExtra int32
	CbWndExtra int32
	HInstance uintptr
	HIcon uintptr
	HCursor uintptr
	HbrBackground uintptr
	LpszMenuName *uint16
	LpszClassName *uint16
	HIconSm uintptr
}
type INITCOMMONCONTROLSEX struct { DwSize uint32; DwICC uint32 }
type LVCOLUMN struct { Mask uint32; Fmt int32; Cx int32; PszText *uint16; CchTextMax int32; ISubItem int32; IImage int32; IOrder int32; CxMin int32; CxDefault int32; CxIdeal int32 }
type LVITEM struct { Mask uint32; IItem int32; ISubItem int32; State uint32; StateMask uint32; PszText *uint16; CchTextMax int32; IImage int32; LParam uintptr; IIndent int32; IGroupId int32; CColumns uint32; PuColumns *uint32; PiColFmt *int32; IGroup int32 }
type NMHDR struct { HwndFrom uintptr; IdFrom uintptr; Code int32 }
type MINMAXINFO struct{ PtReserved, PtMaxSize, PtMaxPosition, PtMinTrackSize, PtMaxTrackSize POINT }
type ACCEL struct { FVirt byte; Key uint16; Cmd uint16 }
type OPENFILENAME struct { LStructSize uint32; HwndOwner uintptr; HInstance uintptr; LpstrFilter *uint16; LpstrCustomFilter *uint16; NMaxCustFilter uint32; NFilterIndex uint32; LpstrFile *uint16; NMaxFile uint32; LpstrFileTitle *uint16; NMaxFileTitle uint32; LpstrInitialDir *uint16; LpstrTitle *uint16; Flags uint32; NFileOffset uint16; NFileExtension uint16; LpstrDefExt *uint16; LCustData uintptr; LpfnHook uintptr; LpTemplateName *uint16; PvReserved uintptr; DwReserved uint32; FlagsEx uint32 }
type BROWSEINFO struct { HwndOwner uintptr; PidlRoot uintptr; PszDisplayName *uint16; LpszTitle *uint16; UlFlags uint32; Lpfn uintptr; LParam uintptr; IImage int32 }
type LOGFONT struct { Height, Width, Escapement, Orientation, Weight int32; Italic, Underline, StrikeOut, CharSet, OutPrecision, ClipPrecision, Quality, PitchAndFamily byte; FaceName [32]uint16 }

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32 = syscall.NewLazyDLL("gdi32.dll")
	shell32 = syscall.NewLazyDLL("shell32.dll")
	comctl32 = syscall.NewLazyDLL("comctl32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	ole32 = syscall.NewLazyDLL("ole32.dll")
	uxtheme = syscall.NewLazyDLL("uxtheme.dll")

	pRegisterClassExW = user32.NewProc("RegisterClassExW")
	pCreateWindowExW = user32.NewProc("CreateWindowExW")
	pDefWindowProcW = user32.NewProc("DefWindowProcW")
	pShowWindow = user32.NewProc("ShowWindow")
	pUpdateWindow = user32.NewProc("UpdateWindow")
	pGetMessageW = user32.NewProc("GetMessageW")
	pTranslateMessage = user32.NewProc("TranslateMessage")
	pDispatchMessageW = user32.NewProc("DispatchMessageW")
	pPostQuitMessage = user32.NewProc("PostQuitMessage")
	pSendMessageW = user32.NewProc("SendMessageW")
	pPostMessageW = user32.NewProc("PostMessageW")
	pMoveWindow = user32.NewProc("MoveWindow")
	pGetClientRect = user32.NewProc("GetClientRect")
	pLoadCursorW = user32.NewProc("LoadCursorW")
	pLoadIconW = user32.NewProc("LoadIconW")
	pLoadImageW = user32.NewProc("LoadImageW")
	pMessageBoxW = user32.NewProc("MessageBoxW")
	pSetWindowTextW = user32.NewProc("SetWindowTextW")
	pGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	pGetWindowTextW = user32.NewProc("GetWindowTextW")
	pEnableWindow = user32.NewProc("EnableWindow")
	pSetTimer = user32.NewProc("SetTimer")
	pKillTimer = user32.NewProc("KillTimer")
	pCreateMenu = user32.NewProc("CreateMenu")
	pCreatePopupMenu = user32.NewProc("CreatePopupMenu")
	pAppendMenuW = user32.NewProc("AppendMenuW")
	pSetMenu = user32.NewProc("SetMenu")
	pTrackPopupMenu = user32.NewProc("TrackPopupMenu")
	pGetCursorPos = user32.NewProc("GetCursorPos")
	pSetProcessDPIAware = user32.NewProc("SetProcessDPIAware")
	pCreateAcceleratorTableW = user32.NewProc("CreateAcceleratorTableW")
	pTranslateAcceleratorW = user32.NewProc("TranslateAcceleratorW")
	pDestroyAcceleratorTable = user32.NewProc("DestroyAcceleratorTable")
	pSetFocus = user32.NewProc("SetFocus")
	pSystemParametersInfoW = user32.NewProc("SystemParametersInfoW")
	pFillRect = user32.NewProc("FillRect")
	pDrawTextW = user32.NewProc("DrawTextW")
	pSetTextColor = gdi32.NewProc("SetTextColor")
	pSetBkMode = gdi32.NewProc("SetBkMode")
	pCreateSolidBrush = gdi32.NewProc("CreateSolidBrush")
	pCreatePen = gdi32.NewProc("CreatePen")
	pSelectObject = gdi32.NewProc("SelectObject")
	pRoundRect = gdi32.NewProc("RoundRect")

	pGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")
	pCreateMutexW = kernel32.NewProc("CreateMutexW")
	pGetLastError = kernel32.NewProc("GetLastError")
	pCloseHandle = kernel32.NewProc("CloseHandle")
	pMultiByteToWideChar = kernel32.NewProc("MultiByteToWideChar")

	pCreateFontW = gdi32.NewProc("CreateFontW")
	pCreateFontIndirectW = gdi32.NewProc("CreateFontIndirectW")
	pDeleteObject = gdi32.NewProc("DeleteObject")

	pDragAcceptFiles = shell32.NewProc("DragAcceptFiles")
	pDragQueryFileW = shell32.NewProc("DragQueryFileW")
	pDragFinish = shell32.NewProc("DragFinish")
	pShellExecuteW = shell32.NewProc("ShellExecuteW")
	pSHBrowseForFolderW = shell32.NewProc("SHBrowseForFolderW")
	pSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")

	pInitCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")
	pGetOpenFileNameW = comdlg32.NewProc("GetOpenFileNameW")
	pGetSaveFileNameW = comdlg32.NewProc("GetSaveFileNameW")
	pCoTaskMemFree = ole32.NewProc("CoTaskMemFree")
	pSetWindowTheme = uxtheme.NewProc("SetWindowTheme")
)

type PdfItem struct { Path string; Pages int; Known bool }
type MergeResult struct { Success bool; Output string; Message string }

var app struct {
	hwnd uintptr
	list, countLabel, titleLabel, subLabel, toolsBtn uintptr
	outputCard, outputTitle uintptr
	saveLabel, saveEdit, browseBtn uintptr
	nameLabel, nameEdit, locateCheck uintptr
	finalLabel, progress, statusLabel, qpdfLabel, mergeBtn uintptr
	buttons []uintptr
	font, fontSmall, fontBold, fontTitle uintptr
	bgBrush, cardBrush uintptr
	menu uintptr
	mu sync.Mutex
	files []PdfItem
	fixedDir string
	manualDir bool
	manualOrder bool
	userEditedName bool
	programmaticName bool
	initialized bool
	qpdf string
	merging bool
	queueDir string
	pageQueue chan string
	mergeResult *MergeResult
}
