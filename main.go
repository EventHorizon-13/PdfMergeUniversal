//go:build windows

package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
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
	CtlType    uint32
	CtlID      uint32
	ItemID     uint32
	ItemAction uint32
	ItemState  uint32
	HwndItem   uintptr
	HDC        uintptr
	RcItem     RECT
	ItemData   uintptr
}
type MSG struct {
	Hwnd           uintptr
	Message        uint32
	WParam, LParam uintptr
	Time           uint32
	Pt             POINT
}
type WNDCLASSEX struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}
type INITCOMMONCONTROLSEX struct {
	DwSize uint32
	DwICC  uint32
}
type LVCOLUMN struct {
	Mask       uint32
	Fmt        int32
	Cx         int32
	PszText    *uint16
	CchTextMax int32
	ISubItem   int32
	IImage     int32
	IOrder     int32
	CxMin      int32
	CxDefault  int32
	CxIdeal    int32
}
type LVITEM struct {
	Mask       uint32
	IItem      int32
	ISubItem   int32
	State      uint32
	StateMask  uint32
	PszText    *uint16
	CchTextMax int32
	IImage     int32
	LParam     uintptr
	IIndent    int32
	IGroupId   int32
	CColumns   uint32
	PuColumns  *uint32
	PiColFmt   *int32
	IGroup     int32
}
type NMHDR struct {
	HwndFrom uintptr
	IdFrom   uintptr
	Code     int32
}
type MINMAXINFO struct{ PtReserved, PtMaxSize, PtMaxPosition, PtMinTrackSize, PtMaxTrackSize POINT }
type ACCEL struct {
	FVirt byte
	Key   uint16
	Cmd   uint16
}
type OPENFILENAME struct {
	LStructSize       uint32
	HwndOwner         uintptr
	HInstance         uintptr
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
	PvReserved        uintptr
	DwReserved        uint32
	FlagsEx           uint32
}
type BROWSEINFO struct {
	HwndOwner      uintptr
	PidlRoot       uintptr
	PszDisplayName *uint16
	LpszTitle      *uint16
	UlFlags        uint32
	Lpfn           uintptr
	LParam         uintptr
	IImage         int32
}
type LOGFONT struct {
	Height, Width, Escapement, Orientation, Weight                                              int32
	Italic, Underline, StrikeOut, CharSet, OutPrecision, ClipPrecision, Quality, PitchAndFamily byte
	FaceName                                                                                    [32]uint16
}

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	gdi32    = syscall.NewLazyDLL("gdi32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")
	comctl32 = syscall.NewLazyDLL("comctl32.dll")
	comdlg32 = syscall.NewLazyDLL("comdlg32.dll")
	ole32    = syscall.NewLazyDLL("ole32.dll")
	uxtheme  = syscall.NewLazyDLL("uxtheme.dll")

	pRegisterClassExW        = user32.NewProc("RegisterClassExW")
	pCreateWindowExW         = user32.NewProc("CreateWindowExW")
	pDefWindowProcW          = user32.NewProc("DefWindowProcW")
	pShowWindow              = user32.NewProc("ShowWindow")
	pUpdateWindow            = user32.NewProc("UpdateWindow")
	pGetMessageW             = user32.NewProc("GetMessageW")
	pTranslateMessage        = user32.NewProc("TranslateMessageW")
	pDispatchMessageW        = user32.NewProc("DispatchMessageW")
	pPostQuitMessage         = user32.NewProc("PostQuitMessage")
	pSendMessageW            = user32.NewProc("SendMessageW")
	pPostMessageW            = user32.NewProc("PostMessageW")
	pMoveWindow              = user32.NewProc("MoveWindow")
	pGetClientRect           = user32.NewProc("GetClientRect")
	pLoadCursorW             = user32.NewProc("LoadCursorW")
	pLoadIconW               = user32.NewProc("LoadIconW")
	pLoadImageW              = user32.NewProc("LoadImageW")
	pMessageBoxW             = user32.NewProc("MessageBoxW")
	pSetWindowTextW          = user32.NewProc("SetWindowTextW")
	pGetWindowTextLengthW    = user32.NewProc("GetWindowTextLengthW")
	pGetWindowTextW          = user32.NewProc("GetWindowTextW")
	pEnableWindow            = user32.NewProc("EnableWindow")
	pSetTimer                = user32.NewProc("SetTimer")
	pKillTimer               = user32.NewProc("KillTimer")
	pCreateMenu              = user32.NewProc("CreateMenu")
	pCreatePopupMenu         = user32.NewProc("CreatePopupMenu")
	pAppendMenuW             = user32.NewProc("AppendMenuW")
	pSetMenu                 = user32.NewProc("SetMenu")
	pTrackPopupMenu          = user32.NewProc("TrackPopupMenu")
	pGetCursorPos            = user32.NewProc("GetCursorPos")
	pSetProcessDPIAware      = user32.NewProc("SetProcessDPIAware")
	pCreateAcceleratorTableW = user32.NewProc("CreateAcceleratorTableW")
	pTranslateAcceleratorW   = user32.NewProc("TranslateAcceleratorW")
	pDestroyAcceleratorTable = user32.NewProc("DestroyAcceleratorTable")
	pSetFocus                = user32.NewProc("SetFocus")
	pSystemParametersInfoW   = user32.NewProc("SystemParametersInfoW")
	pFillRect                = user32.NewProc("FillRect")
	pDrawTextW               = user32.NewProc("DrawTextW")
	pSetTextColor            = gdi32.NewProc("SetTextColor")
	pSetBkMode               = gdi32.NewProc("SetBkMode")
	pCreateSolidBrush        = gdi32.NewProc("CreateSolidBrush")
	pCreatePen               = gdi32.NewProc("CreatePen")
	pSelectObject            = gdi32.NewProc("SelectObject")
	pRoundRect               = gdi32.NewProc("RoundRect")

	pGetModuleHandleW    = kernel32.NewProc("GetModuleHandleW")
	pCreateMutexW        = kernel32.NewProc("CreateMutexW")
	pGetLastError        = kernel32.NewProc("GetLastError")
	pCloseHandle         = kernel32.NewProc("CloseHandle")
	pMultiByteToWideChar = kernel32.NewProc("MultiByteToWideChar")

	pCreateFontW         = gdi32.NewProc("CreateFontW")
	pCreateFontIndirectW = gdi32.NewProc("CreateFontIndirectW")
	pDeleteObject        = gdi32.NewProc("DeleteObject")

	pDragAcceptFiles      = shell32.NewProc("DragAcceptFiles")
	pDragQueryFileW       = shell32.NewProc("DragQueryFileW")
	pDragFinish           = shell32.NewProc("DragFinish")
	pShellExecuteW        = shell32.NewProc("ShellExecuteW")
	pSHBrowseForFolderW   = shell32.NewProc("SHBrowseForFolderW")
	pSHGetPathFromIDListW = shell32.NewProc("SHGetPathFromIDListW")

	pInitCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")
	pGetOpenFileNameW     = comdlg32.NewProc("GetOpenFileNameW")
	pGetSaveFileNameW     = comdlg32.NewProc("GetSaveFileNameW")
	pCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
	pSetWindowTheme       = uxtheme.NewProc("SetWindowTheme")
)

type PdfItem struct {
	Path  string
	Pages int
	Known bool
}

type MergeResult struct {
	Success bool
	Output  string
	Message string
}

var app struct {
	hwnd                                                   uintptr
	list, countLabel, titleLabel, subLabel, toolsBtn       uintptr
	outputCard, outputTitle                                uintptr
	saveLabel, saveEdit, browseBtn                         uintptr
	nameLabel, nameEdit, locateCheck                       uintptr
	finalLabel, progress, statusLabel, qpdfLabel, mergeBtn uintptr
	buttons                                                []uintptr
	font, fontSmall, fontBold, fontTitle                   uintptr
	bgBrush, cardBrush                                     uintptr
	menu                                                   uintptr

	mu               sync.Mutex
	files            []PdfItem
	fixedDir         string
	manualDir        bool
	manualOrder      bool
	userEditedName   bool
	programmaticName bool
	initialized      bool
	qpdf             string
	merging          bool
	queueDir         string
	pageQueue        chan string
	mergeResult      *MergeResult
}

func wstr(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }

func utf16Multi(s string) []uint16 {
	if !strings.HasSuffix(s, "\x00\x00") {
		s += "\x00\x00"
	}
	return utf16.Encode([]rune(s))
}
func loword(v uintptr) uint16 { return uint16(v & 0xffff) }
func hiword(v uintptr) uint16 { return uint16((v >> 16) & 0xffff) }
func boolPtr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func msgbox(title, text string, flags uintptr) int {
	r, _, _ := pMessageBoxW.Call(app.hwnd, uintptr(unsafe.Pointer(wstr(text))), uintptr(unsafe.Pointer(wstr(title))), flags)
	return int(r)
}

func setText(hwnd uintptr, s string) { pSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(wstr(s)))) }
func getText(hwnd uintptr) string {
	n, _, _ := pGetWindowTextLengthW.Call(hwnd)
	if n == 0 {
		return ""
	}
	buf := make([]uint16, n+1)
	pGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), n+1)
	return syscall.UTF16ToString(buf)
}

func send(hwnd uintptr, msg uint32, wp, lp uintptr) uintptr {
	r, _, _ := pSendMessageW.Call(hwnd, uintptr(msg), wp, lp)
	return r
}
func move(hwnd uintptr, x, y, w, h int32) {
	if hwnd != 0 {
		pMoveWindow.Call(hwnd, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 1)
	}
}
func enable(hwnd uintptr, b bool) { pEnableWindow.Call(hwnd, boolPtr(b)) }

func createControl(exStyle uint32, class, text string, style uint32, id int32) uintptr {
	hinst, _, _ := pGetModuleHandleW.Call(0)
	h, _, _ := pCreateWindowExW.Call(uintptr(exStyle), uintptr(unsafe.Pointer(wstr(class))), uintptr(unsafe.Pointer(wstr(text))), uintptr(style|WS_CHILD|WS_VISIBLE), 0, 0, 0, 0, app.hwnd, uintptr(id), hinst, 0)
	if h != 0 && app.font != 0 {
		send(h, WM_SETFONT, app.font, 1)
	}
	if h != 0 {
		pSetWindowTheme.Call(h, uintptr(unsafe.Pointer(wstr("Explorer"))), 0)
	}
	return h
}

func setFont(hwnd, font uintptr) { send(hwnd, WM_SETFONT, font, 1) }

func createFonts() {
	var lf LOGFONT
	ok, _, _ := pSystemParametersInfoW.Call(SPI_GETICONTITLELOGFONT, uintptr(unsafe.Sizeof(lf)), uintptr(unsafe.Pointer(&lf)), 0)
	if ok == 0 {
		face := syscall.StringToUTF16("Segoe UI")
		copy(lf.FaceName[:], face)
		lf.Height, lf.Weight, lf.CharSet, lf.Quality = -12, 400, 1, 5
	}
	app.font, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&lf)))
	small := lf
	app.fontSmall, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&small)))
	bold := lf
	bold.Weight = 600
	app.fontBold, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&bold)))
	title := lf
	title.Height = title.Height * 3 / 2
	title.Weight = 650
	app.fontTitle, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&title)))
}

func rgb(r, g, b byte) uintptr { return uintptr(r) | uintptr(g)<<8 | uintptr(b)<<16 }

func drawButton(di *DRAWITEMSTRUCT) {
	if di == nil {
		return
	}
	if di.HwndItem == app.outputCard {
		br, _, _ := pCreateSolidBrush.Call(rgb(255, 255, 255))
		pen, _, _ := pCreatePen.Call(PS_SOLID, 1, rgb(234, 236, 240))
		oldBr, _, _ := pSelectObject.Call(di.HDC, br)
		oldPen, _, _ := pSelectObject.Call(di.HDC, pen)
		pRoundRect.Call(di.HDC, uintptr(di.RcItem.Left), uintptr(di.RcItem.Top), uintptr(di.RcItem.Right), uintptr(di.RcItem.Bottom), 14, 14)
		pSelectObject.Call(di.HDC, oldBr)
		pSelectObject.Call(di.HDC, oldPen)
		pDeleteObject.Call(br)
		pDeleteObject.Call(pen)
		return
	}
	primary := int(di.CtlID) == ID_MERGE
	disabled := (di.ItemState & ODS_DISABLED) != 0
	pressed := (di.ItemState & ODS_SELECTED) != 0
	bg := rgb(255, 255, 255)
	fg := rgb(52, 64, 84)
	border := rgb(208, 213, 221)
	if primary {
		bg = rgb(37, 99, 235)
		fg = rgb(255, 255, 255)
		border = rgb(37, 99, 235)
		if pressed {
			bg = rgb(29, 78, 216)
			border = bg
		}
	} else if pressed {
		bg = rgb(242, 244, 247)
	}
	if disabled {
		bg = rgb(242, 244, 247)
		fg = rgb(152, 162, 179)
		border = rgb(228, 231, 236)
	}
	br, _, _ := pCreateSolidBrush.Call(bg)
	pen, _, _ := pCreatePen.Call(PS_SOLID, 1, border)
	oldBr, _, _ := pSelectObject.Call(di.HDC, br)
	oldPen, _, _ := pSelectObject.Call(di.HDC, pen)
	pRoundRect.Call(di.HDC, uintptr(di.RcItem.Left), uintptr(di.RcItem.Top), uintptr(di.RcItem.Right), uintptr(di.RcItem.Bottom), 10, 10)
	pSelectObject.Call(di.HDC, oldBr)
	pSelectObject.Call(di.HDC, oldPen)
	pDeleteObject.Call(br)
	pDeleteObject.Call(pen)
	pSetBkMode.Call(di.HDC, TRANSPARENT)
	pSetTextColor.Call(di.HDC, fg)
	font := app.font
	if primary {
		font = app.fontBold
	}
	oldFont, _, _ := pSelectObject.Call(di.HDC, font)
	text := getText(di.HwndItem)
	rr := di.RcItem
	pDrawTextW.Call(di.HDC, uintptr(unsafe.Pointer(wstr(text))), ^uintptr(0), uintptr(unsafe.Pointer(&rr)), DT_CENTER|DT_VCENTER|DT_SINGLELINE)
	pSelectObject.Call(di.HDC, oldFont)
}

func staticColor(hdc, child uintptr) uintptr {
	pSetBkMode.Call(hdc, TRANSPARENT)
	if child == app.countLabel {
		pSetTextColor.Call(hdc, rgb(37, 99, 235))
	} else if child == app.subLabel || child == app.statusLabel || child == app.qpdfLabel || child == app.finalLabel {
		pSetTextColor.Call(hdc, rgb(102, 112, 133))
	} else {
		pSetTextColor.Call(hdc, rgb(16, 24, 40))
	}
	return app.bgBrush
}

func initListColumns() {
	cols := []struct {
		text  string
		width int32
		fmt   int32
	}{{"序号", 55, LVCFMT_LEFT}, {"文件名", 425, LVCFMT_LEFT}, {"页数", 65, LVCFMT_RIGHT}, {"所在文件夹", 400, LVCFMT_LEFT}}
	for i, c := range cols {
		txt := syscall.StringToUTF16(c.text)
		col := LVCOLUMN{Mask: LVCF_TEXT | LVCF_WIDTH | LVCF_FMT, Fmt: c.fmt, Cx: c.width, PszText: &txt[0]}
		send(app.list, LVM_INSERTCOLUMNW, uintptr(i), uintptr(unsafe.Pointer(&col)))
	}
	send(app.list, LVM_SETEXTENDEDLISTVIEWSTYLE, 0, LVS_EX_FULLROWSELECT|LVS_EX_DOUBLEBUFFER)
}

func insertListRow(index int, item PdfItem) {
	vals := []string{strconv.Itoa(index + 1), filepath.Base(item.Path), "?", filepath.Dir(item.Path)}
	if item.Known {
		vals[2] = strconv.Itoa(item.Pages)
	}
	for sub, s := range vals {
		txt := syscall.StringToUTF16(s)
		li := LVITEM{Mask: LVIF_TEXT, IItem: int32(index), ISubItem: int32(sub), PszText: &txt[0]}
		if sub == 0 {
			send(app.list, LVM_INSERTITEMW, 0, uintptr(unsafe.Pointer(&li)))
		} else {
			send(app.list, LVM_SETITEMW, 0, uintptr(unsafe.Pointer(&li)))
		}
	}
}

func refreshList() {
	app.mu.Lock()
	snapshot := append([]PdfItem(nil), app.files...)
	app.mu.Unlock()
	send(app.list, LVM_DELETEALLITEMS, 0, 0)
	total, known := 0, 0
	for i, it := range snapshot {
		insertListRow(i, it)
		if it.Known {
			total += it.Pages
			known++
		}
	}
	if known == len(snapshot) && len(snapshot) > 0 {
		setText(app.countLabel, fmt.Sprintf("%d 个 PDF  ·  %d 页", len(snapshot), total))
	} else if len(snapshot) == 0 {
		setText(app.countLabel, "0 个 PDF  ·  0 页")
	} else {
		setText(app.countLabel, fmt.Sprintf("%d 个 PDF  ·  已识别 %d 页", len(snapshot), total))
	}
	setText(app.statusLabel, fmt.Sprintf("已加载 %d 个 PDF", len(snapshot)))
	updateOutputPreview()
}

func naturalKey(s string) string {
	re := regexp.MustCompile(`\d+`)
	return re.ReplaceAllStringFunc(strings.ToLower(s), func(x string) string { return fmt.Sprintf("%020s", x) })
}

func sortNaturalLocked() {
	sort.SliceStable(app.files, func(i, j int) bool {
		return naturalKey(filepath.Base(app.files[i].Path)) < naturalKey(filepath.Base(app.files[j].Path))
	})
}

func commonAutoName(paths []string) string {
	if len(paths) == 0 {
		return "合并"
	}
	bases := make([]string, 0, len(paths))
	for _, p := range paths {
		bases = append(bases, strings.TrimSuffix(filepath.Base(p), filepath.Ext(p)))
	}
	if len(bases) == 1 {
		s := bases[0]
		s = regexp.MustCompile(`[-_.\s]+\d+$`).ReplaceAllString(s, "")
		s = strings.TrimRight(s, "-_.~、，,；;：:（）()[]{} ")
		if s == "" {
			return "合并"
		}
		return s
	}
	prefix := []rune(bases[0])
	for _, b := range bases[1:] {
		r := []rune(b)
		n := len(prefix)
		if len(r) < n {
			n = len(r)
		}
		k := 0
		for k < n && prefix[k] == r[k] {
			k++
		}
		prefix = prefix[:k]
		if len(prefix) == 0 {
			break
		}
	}
	s := string(prefix)
	if regexp.MustCompile(`\d$`).MatchString(s) {
		s = regexp.MustCompile(`\d+$`).ReplaceAllString(s, "")
	}
	s = strings.TrimRight(s, "-_.~、，,；;：:（）()[]{} ")
	if strings.TrimSpace(s) == "" {
		return "合并"
	}
	return s
}

func maybeAutoName() {
	app.mu.Lock()
	if app.userEditedName {
		app.mu.Unlock()
		return
	}
	paths := make([]string, len(app.files))
	for i, f := range app.files {
		paths[i] = f.Path
	}
	app.mu.Unlock()
	app.programmaticName = true
	setText(app.nameEdit, commonAutoName(paths))
	app.programmaticName = false
	updateOutputPreview()
}

func addPaths(paths []string) {
	var toQueue []string
	app.mu.Lock()
	existing := map[string]bool{}
	for _, f := range app.files {
		existing[strings.ToLower(filepath.Clean(f.Path))] = true
	}
	for _, p := range paths {
		p = filepath.Clean(p)
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if info.IsDir() {
			ents, _ := os.ReadDir(p)
			var nested []string
			for _, e := range ents {
				if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".pdf") {
					nested = append(nested, filepath.Join(p, e.Name()))
				}
			}
			sort.Slice(nested, func(i, j int) bool {
				return naturalKey(filepath.Base(nested[i])) < naturalKey(filepath.Base(nested[j]))
			})
			for _, np := range nested {
				key := strings.ToLower(filepath.Clean(np))
				if existing[key] {
					continue
				}
				existing[key] = true
				app.files = append(app.files, PdfItem{Path: np})
				toQueue = append(toQueue, np)
				if app.fixedDir == "" {
					app.fixedDir = filepath.Dir(np)
				}
			}
			continue
		}
		if !strings.EqualFold(filepath.Ext(p), ".pdf") {
			continue
		}
		key := strings.ToLower(p)
		if existing[key] {
			continue
		}
		existing[key] = true
		app.files = append(app.files, PdfItem{Path: p})
		toQueue = append(toQueue, p)
		if app.fixedDir == "" {
			app.fixedDir = filepath.Dir(p)
		}
	}
	if !app.manualOrder {
		sortNaturalLocked()
	}
	fixed := app.fixedDir
	app.mu.Unlock()
	if fixed != "" && !app.manualDir {
		setText(app.saveEdit, fixed)
	}
	refreshList()
	maybeAutoName()
	for _, p := range toQueue {
		select {
		case app.pageQueue <- p:
		default:
			go func(x string) { app.pageQueue <- x }(p)
		}
	}
}

func selectedIndices() []int {
	var out []int
	idx := -1
	for {
		r := send(app.list, LVM_GETNEXTITEM, uintptr(idx), LVNI_SELECTED)
		ni := int(int32(r))
		if ni < 0 {
			break
		}
		out = append(out, ni)
		idx = ni
	}
	return out
}

func selectIndices(indices []int) {
	count := int(send(app.list, LVM_GETITEMCOUNT, 0, 0))
	for i := 0; i < count; i++ {
		li := LVITEM{StateMask: LVIS_SELECTED | LVIS_FOCUSED, State: 0}
		send(app.list, LVM_SETITEMSTATE, uintptr(i), uintptr(unsafe.Pointer(&li)))
	}
	for _, i := range indices {
		li := LVITEM{StateMask: LVIS_SELECTED | LVIS_FOCUSED, State: LVIS_SELECTED}
		send(app.list, LVM_SETITEMSTATE, uintptr(i), uintptr(unsafe.Pointer(&li)))
	}
	if len(indices) > 0 {
		send(app.list, LVM_ENSUREVISIBLE, uintptr(indices[0]), 0)
	}
}

func deleteSelected() {
	sel := selectedIndices()
	if len(sel) == 0 {
		return
	}
	mark := map[int]bool{}
	for _, i := range sel {
		mark[i] = true
	}
	app.mu.Lock()
	nf := make([]PdfItem, 0, len(app.files)-len(sel))
	for i, f := range app.files {
		if !mark[i] {
			nf = append(nf, f)
		}
	}
	app.files = nf
	app.manualOrder = true
	app.mu.Unlock()
	refreshList()
	maybeAutoName()
}

func moveSelected(delta int) {
	sel := selectedIndices()
	if len(sel) == 0 {
		return
	}
	app.mu.Lock()
	n := len(app.files)
	if delta < 0 {
		selected := map[int]bool{}
		for _, i := range sel {
			selected[i] = true
		}
		for i := 1; i < n; i++ {
			if selected[i] && !selected[i-1] {
				app.files[i-1], app.files[i] = app.files[i], app.files[i-1]
				selected[i-1] = true
				delete(selected, i)
			}
		}
	} else {
		selected := map[int]bool{}
		for _, i := range sel {
			selected[i] = true
		}
		for i := n - 2; i >= 0; i-- {
			if selected[i] && !selected[i+1] {
				app.files[i+1], app.files[i] = app.files[i], app.files[i+1]
				selected[i+1] = true
				delete(selected, i)
			}
		}
	}
	app.manualOrder = true
	app.mu.Unlock()
	refreshList()
	ns := make([]int, 0, len(sel))
	for _, i := range sel {
		j := i + delta
		if j < 0 {
			j = 0
		}
		if j >= n {
			j = n - 1
		}
		ns = append(ns, j)
	}
	selectIndices(ns)
}

func moveTopBottom(top bool) {
	sel := selectedIndices()
	if len(sel) == 0 {
		return
	}
	mark := map[int]bool{}
	for _, i := range sel {
		mark[i] = true
	}
	app.mu.Lock()
	chosen := []PdfItem{}
	rest := []PdfItem{}
	for i, f := range app.files {
		if mark[i] {
			chosen = append(chosen, f)
		} else {
			rest = append(rest, f)
		}
	}
	if top {
		app.files = append(chosen, rest...)
	} else {
		app.files = append(rest, chosen...)
	}
	app.manualOrder = true
	nrest := len(rest)
	nch := len(chosen)
	app.mu.Unlock()
	refreshList()
	idx := []int{}
	if top {
		for i := 0; i < nch; i++ {
			idx = append(idx, i)
		}
	} else {
		for i := 0; i < nch; i++ {
			idx = append(idx, nrest+i)
		}
	}
	selectIndices(idx)
}

func doNaturalSort() {
	app.mu.Lock()
	sortNaturalLocked()
	app.manualOrder = false
	app.mu.Unlock()
	refreshList()
}
func doReverse() {
	app.mu.Lock()
	for i, j := 0, len(app.files)-1; i < j; i, j = i+1, j-1 {
		app.files[i], app.files[j] = app.files[j], app.files[i]
	}
	app.manualOrder = true
	app.mu.Unlock()
	refreshList()
}
func clearList() {
	app.mu.Lock()
	app.files = nil
	app.fixedDir = ""
	app.manualDir = false
	app.manualOrder = false
	app.userEditedName = false
	app.mu.Unlock()
	setText(app.saveEdit, "")
	app.programmaticName = true
	setText(app.nameEdit, "合并")
	app.programmaticName = false
	refreshList()
}

func sanitizeName(s string) string {
	s = strings.TrimSpace(s)
	if strings.EqualFold(filepath.Ext(s), ".pdf") {
		s = strings.TrimSuffix(s, filepath.Ext(s))
	}
	bad := regexp.MustCompile(`[<>:"/\\|?*]`)
	s = bad.ReplaceAllString(s, "_")
	s = strings.TrimRight(s, ". ")
	if s == "" {
		s = "合并"
	}
	return s
}

func outputPath() string {
	dir := strings.TrimSpace(getText(app.saveEdit))
	name := sanitizeName(getText(app.nameEdit))
	if dir == "" {
		return ""
	}
	p := filepath.Join(dir, name+".pdf")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return p
	}
	for i := 2; i < 10000; i++ {
		q := filepath.Join(dir, fmt.Sprintf("%s (%d).pdf", name, i))
		if _, err := os.Stat(q); os.IsNotExist(err) {
			return q
		}
	}
	return p
}

func updateOutputPreview() {
	p := outputPath()
	if p == "" {
		setText(app.finalLabel, "最终输出：请先添加 PDF")
		return
	}
	if strings.Contains(filepath.Base(p), " (") {
		setText(app.finalLabel, "最终输出："+p+"   （原文件已存在，将自动改名）")
	} else {
		setText(app.finalLabel, "最终输出："+p)
	}
}

func findQpdf() string {
	if app.qpdf != "" {
		if _, e := os.Stat(app.qpdf); e == nil {
			return app.qpdf
		}
	}
	if p := loadQpdfConfig(); p != "" {
		if _, e := os.Stat(p); e == nil {
			return p
		}
	}
	exe, _ := os.Executable()
	base := filepath.Dir(exe)
	cands := []string{filepath.Join(base, "qpdf.exe"), filepath.Join(base, "qpdf", "bin", "qpdf.exe")}
	for _, env := range []string{"ProgramW6432", "ProgramFiles", "ProgramFiles(x86)"} {
		if v := os.Getenv(env); v != "" {
			matches, _ := filepath.Glob(filepath.Join(v, "qpdf*", "bin", "qpdf.exe"))
			sort.Sort(sort.Reverse(sort.StringSlice(matches)))
			cands = append(cands, matches...)
		}
	}
	if lp, e := exec.LookPath("qpdf.exe"); e == nil {
		cands = append(cands, lp)
	}
	for _, p := range cands {
		if p != "" {
			if _, e := os.Stat(p); e == nil {
				return p
			}
		}
	}
	return ""
}

func configDir() string {
	d := os.Getenv("APPDATA")
	if d == "" {
		d = os.TempDir()
	}
	return filepath.Join(d, "PdfMergeUniversal")
}
func loadQpdfConfig() string {
	b, e := os.ReadFile(filepath.Join(configDir(), "qpdf.txt"))
	if e != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}
func saveQpdfConfig(p string) {
	os.MkdirAll(configDir(), 0755)
	_ = os.WriteFile(filepath.Join(configDir(), "qpdf.txt"), []byte(p), 0644)
}
func updateQpdfLabel() {
	app.qpdf = findQpdf()
	if app.qpdf == "" {
		setText(app.qpdfLabel, "qpdf：未配置（右上角 ··· → 配置 qpdf）")
	} else {
		setText(app.qpdfLabel, "qpdf："+app.qpdf)
	}
}

func pageCount(path string) (int, bool) {
	q := app.qpdf
	if q == "" {
		return 0, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	out, _ := exec.CommandContext(ctx, q, "--warning-exit-0", "--show-npages", path).CombinedOutput()
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if n, e := strconv.Atoi(s); e == nil && n >= 0 {
			return n, true
		}
	}
	return 0, false
}

func startPageWorkers() {
	app.pageQueue = make(chan string, 256)
	for k := 0; k < 2; k++ {
		go func() {
			for p := range app.pageQueue {
				n, ok := pageCount(p)
				app.mu.Lock()
				for i := range app.files {
					if strings.EqualFold(app.files[i].Path, p) {
						app.files[i].Pages = n
						app.files[i].Known = ok
						break
					}
				}
				app.mu.Unlock()
				pPostMessageW.Call(app.hwnd, MSG_REFRESH, 0, 0)
			}
		}()
	}
}

func chooseOpenPDFs() []string {
	buf := make([]uint16, 65536)
	filter := utf16Multi("PDF 文件 (*.pdf)\x00*.pdf\x00所有文件 (*.*)\x00*.*")
	title := syscall.StringToUTF16("添加 PDF")
	ofn := OPENFILENAME{LStructSize: uint32(unsafe.Sizeof(OPENFILENAME{})), HwndOwner: app.hwnd, LpstrFilter: &filter[0], LpstrFile: &buf[0], NMaxFile: uint32(len(buf)), LpstrTitle: &title[0], Flags: OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_PATHMUSTEXIST | OFN_ALLOWMULTISELECT}
	r, _, _ := pGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	runtime.KeepAlive(filter)
	runtime.KeepAlive(title)
	runtime.KeepAlive(buf)
	if r == 0 {
		return nil
	}
	parts := splitMultiSelect(buf)
	if len(parts) == 1 {
		return parts
	}
	dir := parts[0]
	out := []string{}
	for _, f := range parts[1:] {
		out = append(out, filepath.Join(dir, f))
	}
	return out
}

func splitMultiSelect(buf []uint16) []string {
	var out []string
	start := 0
	for i := 0; i < len(buf); i++ {
		if buf[i] == 0 {
			if i == start {
				break
			}
			out = append(out, syscall.UTF16ToString(buf[start:i]))
			start = i + 1
		}
	}
	return out
}

func chooseFile(title, filterStr, defExt string, save bool) string {
	buf := make([]uint16, 4096)
	filterW := utf16Multi(filterStr)
	titleW := syscall.StringToUTF16(title)
	defExtW := syscall.StringToUTF16(defExt)
	flags := uint32(OFN_EXPLORER | OFN_PATHMUSTEXIST)
	if save {
		flags |= OFN_OVERWRITEPROMPT
	} else {
		flags |= OFN_FILEMUSTEXIST
	}
	ofn := OPENFILENAME{
		LStructSize: uint32(unsafe.Sizeof(OPENFILENAME{})),
		HwndOwner:   app.hwnd,
		LpstrFilter: &filterW[0],
		LpstrFile:   &buf[0],
		NMaxFile:    uint32(len(buf)),
		LpstrTitle:  &titleW[0],
		Flags:       flags,
		LpstrDefExt: &defExtW[0],
	}
	var r uintptr
	if save {
		r, _, _ = pGetSaveFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	} else {
		r, _, _ = pGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	}
	runtime.KeepAlive(filterW)
	runtime.KeepAlive(titleW)
	runtime.KeepAlive(defExtW)
	runtime.KeepAlive(buf)
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func chooseFolder() string {
	display := make([]uint16, 260)
	titleW := syscall.StringToUTF16("选择合并后 PDF 的保存文件夹")
	bi := BROWSEINFO{HwndOwner: app.hwnd, PszDisplayName: &display[0], LpszTitle: &titleW[0], UlFlags: BIF_RETURNONLYFSDIRS | BIF_NEWDIALOGSTYLE}
	pidl, _, _ := pSHBrowseForFolderW.Call(uintptr(unsafe.Pointer(&bi)))
	runtime.KeepAlive(titleW)
	runtime.KeepAlive(display)
	if pidl == 0 {
		return ""
	}
	defer pCoTaskMemFree.Call(pidl)
	buf := make([]uint16, 32768)
	r, _, _ := pSHGetPathFromIDListW.Call(pidl, uintptr(unsafe.Pointer(&buf[0])))
	if r == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf)
}

func exportCSV() {
	p := chooseFile("导出顺序 CSV", "CSV 文件 (*.csv)\x00*.csv\x00所有文件 (*.*)\x00*.*", "csv", true)
	if p == "" {
		return
	}
	if !strings.HasSuffix(strings.ToLower(p), ".csv") {
		p += ".csv"
	}
	f, e := os.Create(p)
	if e != nil {
		msgbox("导出失败", e.Error(), MB_OK|MB_ICONERROR)
		return
	}
	defer f.Close()
	f.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(f)
	_ = w.Write([]string{"顺序", "文件名", "完整路径"})
	app.mu.Lock()
	for i, it := range app.files {
		_ = w.Write([]string{strconv.Itoa(i + 1), filepath.Base(it.Path), it.Path})
	}
	app.mu.Unlock()
	w.Flush()
	if e = w.Error(); e != nil {
		msgbox("导出失败", e.Error(), MB_OK|MB_ICONERROR)
		return
	}
	setText(app.statusLabel, "已导出顺序 CSV："+p)
}

func decodeCSVBytes(b []byte) ([]byte, error) {
	b = bytes.TrimPrefix(b, []byte{0xEF, 0xBB, 0xBF})
	if len(b) == 0 || bytes.Contains(b, []byte("文件名")) || bytes.Contains(b, []byte("完整路径")) {
		return b, nil
	}

	const cpGBK = 936
	n, _, _ := pMultiByteToWideChar.Call(cpGBK, 0, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), 0, 0)
	if n == 0 {
		return nil, fmt.Errorf("CSV 既不是有效 UTF-8，也无法按 GBK/ANSI 解码")
	}
	wide := make([]uint16, int(n))
	n2, _, _ := pMultiByteToWideChar.Call(cpGBK, 0, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(unsafe.Pointer(&wide[0])), n)
	if n2 == 0 {
		return nil, fmt.Errorf("GBK/ANSI CSV 解码失败")
	}
	return []byte(string(utf16.Decode(wide[:int(n2)]))), nil
}

func importCSV() {
	p := chooseFile("导入顺序 CSV", "CSV 文件 (*.csv)\x00*.csv\x00所有文件 (*.*)\x00*.*", "csv", false)
	if p == "" {
		return
	}
	b, e := os.ReadFile(p)
	if e != nil {
		msgbox("导入失败", e.Error(), MB_OK|MB_ICONERROR)
		return
	}
	b, e = decodeCSVBytes(b)
	if e != nil {
		msgbox("CSV 编码不支持", e.Error(), MB_OK|MB_ICONWARNING)
		return
	}
	r := csv.NewReader(bytes.NewReader(b))
	rows, e := r.ReadAll()
	if e != nil || len(rows) < 2 {
		msgbox("CSV 格式不正确", "请先使用“导出顺序 CSV”生成模板。", MB_OK|MB_ICONWARNING)
		return
	}
	header := rows[0]
	idxName, idxPath := -1, -1
	for i, h := range header {
		switch strings.TrimSpace(strings.TrimPrefix(h, "\uFEFF")) {
		case "文件名":
			idxName = i
		case "完整路径":
			idxPath = i
		}
	}
	if idxName < 0 && idxPath < 0 {
		msgbox("CSV 格式不正确", "CSV 至少需要“文件名”或“完整路径”列。", MB_OK|MB_ICONWARNING)
		return
	}
	app.mu.Lock()
	current := append([]PdfItem(nil), app.files...)
	app.mu.Unlock()
	used := map[int]bool{}
	ordered := []PdfItem{}
	for _, row := range rows[1:] {
		var name, path string
		if idxName >= 0 && idxName < len(row) {
			name = strings.TrimSpace(row[idxName])
		}
		if idxPath >= 0 && idxPath < len(row) {
			path = strings.TrimSpace(row[idxPath])
		}
		for i, it := range current {
			if used[i] {
				continue
			}
			match := false
			if path != "" && strings.EqualFold(filepath.Clean(path), filepath.Clean(it.Path)) {
				match = true
			} else if path == "" && name != "" && strings.EqualFold(name, filepath.Base(it.Path)) {
				match = true
			}
			if match {
				ordered = append(ordered, it)
				used[i] = true
				break
			}
		}
	}
	for i, it := range current {
		if !used[i] {
			ordered = append(ordered, it)
		}
	}
	app.mu.Lock()
	app.files = ordered
	app.manualOrder = true
	app.mu.Unlock()
	refreshList()
	setText(app.statusLabel, "已按 CSV 调整顺序")
}

func launchPath(path string) {
	pShellExecuteW.Call(0, uintptr(unsafe.Pointer(wstr("open"))), uintptr(unsafe.Pointer(wstr(path))), 0, 0, SW_SHOWNORMAL)
}
func revealPath(path string) { exec.Command("explorer.exe", "/select,"+path).Start() }

func setBusy(b bool) {
	app.merging = b
	for _, h := range app.buttons {
		enable(h, !b)
	}
	enable(app.saveEdit, !b)
	enable(app.nameEdit, !b)
	enable(app.browseBtn, !b)
	enable(app.locateCheck, !b)
	enable(app.mergeBtn, !b)
	if b {
		send(app.progress, PBM_SETRANGE32, 0, 100)
		send(app.progress, PBM_SETPOS, 35, 0)
	} else {
		send(app.progress, PBM_SETPOS, 0, 0)
	}
}

func startMerge() {
	if app.merging {
		return
	}
	app.mu.Lock()
	snapshot := append([]PdfItem(nil), app.files...)
	app.mu.Unlock()
	if len(snapshot) < 2 {
		msgbox("无法合并", "请至少添加两个 PDF 文件。", MB_OK|MB_ICONWARNING)
		return
	}
	if app.qpdf == "" {
		msgbox("未配置 qpdf", "当前没有找到 qpdf.exe。\n请点击右上角“···”，选择“配置 qpdf...”后指定 qpdf.exe。", MB_OK|MB_ICONWARNING)
		return
	}
	out := outputPath()
	if out == "" {
		msgbox("无法合并", "保存路径无效。", MB_OK|MB_ICONWARNING)
		return
	}
	if e := os.MkdirAll(filepath.Dir(out), 0755); e != nil {
		msgbox("无法合并", e.Error(), MB_OK|MB_ICONERROR)
		return
	}
	total := 0
	known := true
	for _, it := range snapshot {
		if !it.Known {
			known = false
		} else {
			total += it.Pages
		}
	}
	setBusy(true)
	if known {
		setText(app.statusLabel, fmt.Sprintf("正在合并 %d 个 PDF，共 %d 页……", len(snapshot), total))
	} else {
		setText(app.statusLabel, fmt.Sprintf("正在合并 %d 个 PDF……", len(snapshot)))
	}
	go func() {
		args := []string{"--warning-exit-0", "--empty", "--pages"}
		for _, it := range snapshot {
			args = append(args, it.Path)
		}
		args = append(args, "--", out)
		cmd := exec.Command(app.qpdf, args...)
		buf, err := cmd.CombinedOutput()
		success := false
		msg := strings.TrimSpace(string(buf))
		if st, e := os.Stat(out); e == nil && st.Size() > 0 {
			success = true
		}
		if err != nil && !success {
			if msg == "" {
				msg = err.Error()
			}
		}
		if success && known {
			if n, ok := pageCount(out); ok && n != total {
				success = false
				msg = fmt.Sprintf("输出页数校验失败：预计 %d 页，实际 %d 页。", total, n)
			}
		}
		app.mu.Lock()
		app.mergeResult = &MergeResult{Success: success, Output: out, Message: msg}
		app.mu.Unlock()
		pPostMessageW.Call(app.hwnd, MSG_MERGE_DONE, 0, 0)
	}()
}

func handleMergeDone() {
	app.mu.Lock()
	r := app.mergeResult
	app.mergeResult = nil
	app.mu.Unlock()
	setBusy(false)
	if r == nil {
		return
	}
	if r.Success {
		setText(app.statusLabel, "合并完成："+r.Output)
		setText(app.progress, "")
		if send(app.locateCheck, 0x00F0, 0, 0) != 0 {
			revealPath(r.Output)
		}
		msgbox("合并完成", "PDF 已成功生成：\n"+r.Output, MB_OK|MB_ICONINFORMATION)
	} else {
		setText(app.statusLabel, "合并失败")
		m := r.Message
		if m == "" {
			m = "qpdf 未生成有效的输出文件。"
		}
		msgbox("合并失败", m, MB_OK|MB_ICONERROR)
	}
	updateOutputPreview()
}

func installContextMenu() {
	exe, e := os.Executable()
	if e != nil {
		return
	}
	key := `HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\PdfMergeUniversal`
	cmdkey := key + `\command`
	cmd := fmt.Sprintf(`"%s" "%%1"`, exe)
	commands := [][]string{{"add", key, "/ve", "/d", "合并 PDF（可调整顺序）", "/f"}, {"add", key, "/v", "MultiSelectModel", "/t", "REG_SZ", "/d", "Player", "/f"}, {"add", key, "/v", "Icon", "/t", "REG_SZ", "/d", exe + ",0", "/f"}, {"add", cmdkey, "/ve", "/d", cmd, "/f"}}
	for _, a := range commands {
		if out, err := exec.Command("reg.exe", a...).CombinedOutput(); err != nil {
			msgbox("安装右键失败", string(out), MB_OK|MB_ICONERROR)
			return
		}
	}
	msgbox("安装完成", "PDF 右键菜单已注册。\nWindows 11 中可能显示在“显示更多选项”里。", MB_OK|MB_ICONINFORMATION)
}
func uninstallContextMenu() {
	key := `HKCU\Software\Classes\SystemFileAssociations\.pdf\shell\PdfMergeUniversal`
	_ = exec.Command("reg.exe", "delete", key, "/f").Run()
	msgbox("卸载完成", "PDF 右键菜单已移除。", MB_OK|MB_ICONINFORMATION)
}
func configureQpdf() {
	p := chooseFile("选择 qpdf.exe", "qpdf.exe\x00qpdf.exe\x00可执行文件 (*.exe)\x00*.exe", "exe", false)
	if p == "" {
		return
	}
	saveQpdfConfig(p)
	app.qpdf = p
	updateQpdfLabel()
	app.mu.Lock()
	unknown := []string{}
	for _, it := range app.files {
		if !it.Known {
			unknown = append(unknown, it.Path)
		}
	}
	app.mu.Unlock()
	for _, x := range unknown {
		app.pageQueue <- x
	}
}

func buildMenu() {}

func showToolsMenu() {
	m, _, _ := pCreatePopupMenu.Call()
	pAppendMenuW.Call(m, MF_STRING, ID_MENU_INSTALL, uintptr(unsafe.Pointer(wstr("安装 PDF 右键菜单"))))
	pAppendMenuW.Call(m, MF_STRING, ID_MENU_UNINSTALL, uintptr(unsafe.Pointer(wstr("卸载 PDF 右键菜单"))))
	pAppendMenuW.Call(m, MF_SEPARATOR, 0, 0)
	pAppendMenuW.Call(m, MF_STRING, ID_MENU_QPDF, uintptr(unsafe.Pointer(wstr("配置 qpdf..."))))
	pAppendMenuW.Call(m, MF_SEPARATOR, 0, 0)
	pAppendMenuW.Call(m, MF_STRING, ID_MENU_EXIT, uintptr(unsafe.Pointer(wstr("退出"))))
	var pt POINT
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	cmd, _, _ := pTrackPopupMenu.Call(m, TPM_RIGHTBUTTON|TPM_RETURNCMD, uintptr(pt.X), uintptr(pt.Y), 0, app.hwnd, 0)
	if cmd != 0 {
		handleCommand(uint16(cmd), 0)
	}
}

func layout() {
	var rc RECT
	pGetClientRect.Call(app.hwnd, uintptr(unsafe.Pointer(&rc)))
	w := rc.Right - rc.Left
	h := rc.Bottom - rc.Top
	margin := int32(24)
	gap := int32(16)
	rightW := int32(142)
	headerTop := int32(20)
	headerH := int32(72)
	bottomH := int32(176)
	if w < 760 || h < 560 {
		return
	}
	move(app.titleLabel, margin, headerTop, w-2*margin-250, 38)
	move(app.countLabel, w-margin-240, headerTop+4, 190, 30)
	move(app.toolsBtn, w-margin-40, headerTop+2, 40, 34)
	move(app.subLabel, margin, headerTop+42, w-2*margin, 24)
	contentTop := headerTop + headerH
	cardTop := h - margin - bottomH
	minListH := int32(250)
	if cardTop < contentTop+minListH+14 {
		cardTop = contentTop + minListH + 14
	}
	listBottom := cardTop - 14
	listH := listBottom - contentTop
	listW := w - 2*margin - rightW - gap
	if listW < 480 {
		listW = 480
	}
	move(app.list, margin, contentTop, listW, listH)
	bx := margin + listW + gap
	by := contentTop
	bh := int32(34)
	buttonCount := int32(len(app.buttons))
	bg := int32(8)
	needed := buttonCount*bh + (buttonCount-1)*bg
	if needed > listH && buttonCount > 1 {
		bg = (listH - buttonCount*bh) / (buttonCount - 1)
		if bg < 2 {
			bg = 2
		}
	}
	for _, b := range app.buttons {
		move(b, bx, by, rightW, bh)
		by += bh + bg
	}
	cardX := margin
	cardW := w - 2*margin
	innerX := cardX + 12
	innerW := cardW - 24
	move(app.outputTitle, innerX, cardTop+4, 160, 24)
	labelW := int32(78)
	browseW := int32(82)
	row1 := cardTop + 34
	move(app.saveLabel, innerX, row1+5, labelW, 24)
	move(app.saveEdit, innerX+labelW, row1, innerW-labelW-browseW-gap, 32)
	move(app.browseBtn, cardX+cardW-12-browseW, row1, browseW, 32)
	row2 := row1 + 40
	checkW := int32(180)
	move(app.nameLabel, innerX, row2+5, labelW, 24)
	nameW := innerW - labelW - checkW - gap
	if nameW > 420 {
		nameW = 420
	}
	if nameW < 220 {
		nameW = 220
	}
	move(app.nameEdit, innerX+labelW, row2, nameW, 32)
	move(app.locateCheck, cardX+cardW-12-checkW, row2+4, checkW, 26)
	mergeW := int32(148)
	row3 := row2 + 39
	move(app.finalLabel, innerX, row3+2, innerW-mergeW-gap, 24)
	move(app.mergeBtn, cardX+cardW-12-mergeW, row3-5, mergeW, 42)
	row4 := row3 + 30
	leftW := innerW - mergeW - gap
	move(app.progress, innerX, row4, leftW, 6)
	move(app.statusLabel, innerX, row4+12, leftW/2, 20)
	move(app.qpdfLabel, innerX+leftW/2, row4+12, leftW/2, 20)
}

func createUI() {
	createFonts()
	app.titleLabel = createControl(0, "STATIC", "PDF 合并工具", 0, 0)
	setFont(app.titleLabel, app.fontTitle)
	app.countLabel = createControl(0, "STATIC", "0 个 PDF · 0 页", SS_CENTER|SS_CENTERIMAGE, 0)
	setFont(app.countLabel, app.fontBold)
	app.toolsBtn = createControl(0, "BUTTON", "···", BS_OWNERDRAW|WS_TABSTOP, ID_TOOLS)
	app.subLabel = createControl(0, "STATIC", "拖入 PDF 或文件夹，调整顺序后按 Enter 即可合并", 0, 0)
	setFont(app.subLabel, app.font)
	app.list = createControl(WS_EX_CLIENTEDGE, "SysListView32", "", LVS_REPORT|LVS_SHOWSELALWAYS|WS_TABSTOP|WS_VSCROLL|WS_HSCROLL, ID_LIST)
	initListColumns()
	labels := []struct {
		id   int32
		text string
	}{{ID_ADD, "＋  添加 PDF"}, {ID_DELETE, "删除选中"}, {ID_UP, "↑  上移"}, {ID_DOWN, "↓  下移"}, {ID_NATSORT, "自然排序"}, {ID_REVERSE, "反转顺序"}, {ID_EXPORT_CSV, "导出顺序 CSV"}, {ID_IMPORT_CSV, "导入顺序 CSV"}, {ID_CLEAR, "清空列表"}}
	for _, x := range labels {
		h := createControl(0, "BUTTON", x.text, BS_OWNERDRAW|WS_TABSTOP, x.id)
		app.buttons = append(app.buttons, h)
	}
	app.outputCard = 0
	app.outputTitle = createControl(0, "STATIC", "输出设置", 0, 0)
	setFont(app.outputTitle, app.fontBold)
	app.saveLabel = createControl(0, "STATIC", "保存位置", 0, 0)
	app.saveEdit = createControl(WS_EX_CLIENTEDGE, "EDIT", "", ES_AUTOHSCROLL|WS_TABSTOP, 0)
	app.browseBtn = createControl(0, "BUTTON", "浏览", BS_OWNERDRAW|WS_TABSTOP, ID_BROWSE_DIR)
	app.nameLabel = createControl(0, "STATIC", "文件名", 0, 0)
	app.nameEdit = createControl(WS_EX_CLIENTEDGE, "EDIT", "合并", ES_AUTOHSCROLL|WS_TABSTOP, ID_NAME)
	app.locateCheck = createControl(0, "BUTTON", "完成后定位文件", BS_AUTOCHECKBOX|WS_TABSTOP, ID_LOCATE)
	send(app.locateCheck, 0x00F1, 0, 0)
	app.finalLabel = createControl(0, "STATIC", "最终输出：请先添加 PDF", 0, 0)
	setFont(app.finalLabel, app.fontSmall)
	app.progress = createControl(0, "msctls_progress32", "", PBS_SMOOTH, 0)
	app.statusLabel = createControl(0, "STATIC", "已加载 0 个 PDF", 0, 0)
	setFont(app.statusLabel, app.fontSmall)
	app.qpdfLabel = createControl(0, "STATIC", "qpdf：检测中...", 0, 0)
	setFont(app.qpdfLabel, app.fontSmall)
	app.mergeBtn = createControl(0, "BUTTON", "合并 PDF", BS_OWNERDRAW|BS_DEFPUSHBUTTON|WS_TABSTOP, ID_MERGE)
	setFont(app.mergeBtn, app.fontBold)
	buildMenu()
	pDragAcceptFiles.Call(app.hwnd, 1)
	pSetTimer.Call(app.hwnd, TIMER_QUEUE, 300, 0)
	app.initialized = true
	updateQpdfLabel()
	layout()
	startPageWorkers()
}

func handleDrop(hdrop uintptr) {
	count, _, _ := pDragQueryFileW.Call(hdrop, 0xFFFFFFFF, 0, 0)
	paths := make([]string, 0, count)
	for i := uintptr(0); i < count; i++ {
		n, _, _ := pDragQueryFileW.Call(hdrop, i, 0, 0)
		buf := make([]uint16, n+1)
		pDragQueryFileW.Call(hdrop, i, uintptr(unsafe.Pointer(&buf[0])), n+1)
		paths = append(paths, syscall.UTF16ToString(buf))
	}
	pDragFinish.Call(hdrop)
	addPaths(paths)
}

func showListContext() {
	sel := selectedIndices()
	if len(sel) == 0 {
		return
	}
	m, _, _ := pCreatePopupMenu.Call()
	items := []struct {
		id   uintptr
		text string
	}{{ID_CTX_OPEN, "打开文件"}, {ID_CTX_FOLDER, "打开所在文件夹"}, {0, ""}, {ID_CTX_TOP, "移到最前"}, {ID_CTX_BOTTOM, "移到最后"}, {ID_CTX_UP, "上移"}, {ID_CTX_DOWN, "下移"}, {0, ""}, {ID_CTX_DELETE, "删除"}}
	for _, it := range items {
		if it.id == 0 {
			pAppendMenuW.Call(m, MF_SEPARATOR, 0, 0)
		} else {
			pAppendMenuW.Call(m, MF_STRING, it.id, uintptr(unsafe.Pointer(wstr(it.text))))
		}
	}
	var pt POINT
	pGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	cmd, _, _ := pTrackPopupMenu.Call(m, TPM_RIGHTBUTTON|TPM_RETURNCMD, uintptr(pt.X), uintptr(pt.Y), 0, app.hwnd, 0)
	switch cmd {
	case ID_CTX_OPEN:
		app.mu.Lock()
		p := app.files[sel[0]].Path
		app.mu.Unlock()
		launchPath(p)
	case ID_CTX_FOLDER:
		app.mu.Lock()
		p := app.files[sel[0]].Path
		app.mu.Unlock()
		revealPath(p)
	case ID_CTX_TOP:
		moveTopBottom(true)
	case ID_CTX_BOTTOM:
		moveTopBottom(false)
	case ID_CTX_UP:
		moveSelected(-1)
	case ID_CTX_DOWN:
		moveSelected(1)
	case ID_CTX_DELETE:
		deleteSelected()
	}
}

func handleCommand(id uint16, code uint16) {
	if id == ID_NAME && code == EN_CHANGE {
		if app.initialized && !app.programmaticName {
			app.userEditedName = true
		}
		updateOutputPreview()
		return
	}
	switch id {
	case ID_ADD, ID_MENU_ADD:
		addPaths(chooseOpenPDFs())
	case ID_TOOLS:
		showToolsMenu()
	case ID_DELETE:
		deleteSelected()
	case ID_UP:
		moveSelected(-1)
	case ID_DOWN:
		moveSelected(1)
	case ID_NATSORT:
		doNaturalSort()
	case ID_REVERSE:
		doReverse()
	case ID_EXPORT_CSV:
		exportCSV()
	case ID_IMPORT_CSV:
		importCSV()
	case ID_CLEAR, ID_MENU_CLEAR:
		clearList()
	case ID_BROWSE_DIR:
		if p := chooseFolder(); p != "" {
			app.manualDir = true
			setText(app.saveEdit, p)
			updateOutputPreview()
		}
	case ID_MERGE:
		startMerge()
	case ID_SELECT_ALL:
		count := int(send(app.list, LVM_GETITEMCOUNT, 0, 0))
		idx := make([]int, count)
		for i := 0; i < count; i++ {
			idx[i] = i
		}
		selectIndices(idx)
	case ID_MENU_INSTALL:
		installContextMenu()
	case ID_MENU_UNINSTALL:
		uninstallContextMenu()
	case ID_MENU_QPDF:
		configureQpdf()
	case ID_MENU_EXIT:
		pPostQuitMessage.Call(0)
	}
}

func scanQueue() {
	ents, _ := os.ReadDir(app.queueDir)
	var paths []string
	for _, e := range ents {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".txt") {
			continue
		}
		p := filepath.Join(app.queueDir, e.Name())
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		_ = os.Remove(p)
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
			if line != "" {
				paths = append(paths, line)
			}
		}
	}
	if len(paths) > 0 {
		addPaths(paths)
	}
}

func queueSecondary(args []string) {
	if len(args) == 0 {
		return
	}
	os.MkdirAll(app.queueDir, 0755)
	name := fmt.Sprintf("%d-%d.txt", os.Getpid(), time.Now().UnixNano())
	var b strings.Builder
	for _, a := range args {
		b.WriteString(a)
		b.WriteByte('\n')
	}
	_ = os.WriteFile(filepath.Join(app.queueDir, name), []byte(b.String()), 0644)
}

func safeHandleCommand(id, code uint16) {
	defer func() {
		if r := recover(); r != nil {
			msgbox("操作失败", fmt.Sprintf("发生内部错误：%v", r), MB_OK|MB_ICONERROR)
		}
	}()
	handleCommand(id, code)
}

func wndProc(hwnd uintptr, msg uint32, wp, lp uintptr) uintptr {
	switch msg {
	case WM_CREATE:
		app.hwnd = hwnd
		createUI()
		return 0
	case WM_SIZE:
		layout()
		return 0
	case WM_GETMINMAXINFO:
		m := (*MINMAXINFO)(unsafe.Pointer(lp))
		m.PtMinTrackSize.X = 940
		m.PtMinTrackSize.Y = 700
		return 0
	case WM_COMMAND:
		safeHandleCommand(loword(wp), hiword(wp))
		return 0
	case WM_DRAWITEM:
		drawButton((*DRAWITEMSTRUCT)(unsafe.Pointer(lp)))
		return 1
	case WM_CTLCOLORSTATIC:
		return staticColor(wp, lp)
	case WM_NOTIFY:
		hdr := (*NMHDR)(unsafe.Pointer(lp))
		if hdr.IdFrom == ID_LIST {
			if hdr.Code == NM_DBLCLK {
				sel := selectedIndices()
				if len(sel) > 0 {
					app.mu.Lock()
					p := app.files[sel[0]].Path
					app.mu.Unlock()
					launchPath(p)
				}
			} else if hdr.Code == NM_RCLICK {
				showListContext()
			}
		}
		return 0
	case WM_DROPFILES:
		handleDrop(wp)
		return 0
	case WM_TIMER:
		if wp == TIMER_QUEUE {
			scanQueue()
		}
		return 0
	case MSG_REFRESH:
		refreshList()
		return 0
	case MSG_MERGE_DONE:
		handleMergeDone()
		return 0
	case WM_DESTROY:
		pKillTimer.Call(hwnd, TIMER_QUEUE)
		pPostQuitMessage.Call(0)
		return 0
	}
	r, _, _ := pDefWindowProcW.Call(hwnd, uintptr(msg), wp, lp)
	return r
}

func createAccelerators() uintptr {
	a := []ACCEL{{FVIRTKEY, VK_RETURN, ID_MERGE}, {FVIRTKEY | FCONTROL, VK_RETURN, ID_MERGE}, {FVIRTKEY | FCONTROL, 'O', ID_ADD}, {FVIRTKEY | FCONTROL, 'A', ID_SELECT_ALL}, {FVIRTKEY, VK_DELETE, ID_DELETE}, {FVIRTKEY | FALT, VK_UP, ID_UP}, {FVIRTKEY | FALT, VK_DOWN, ID_DOWN}}
	h, _, _ := pCreateAcceleratorTableW.Call(uintptr(unsafe.Pointer(&a[0])), uintptr(len(a)))
	return h
}

func main() {
	runtime.LockOSThread()
	pSetProcessDPIAware.Call()
	icc := INITCOMMONCONTROLSEX{DwSize: uint32(unsafe.Sizeof(INITCOMMONCONTROLSEX{})), DwICC: 0x00000001 | 0x00000004 | 0x00004000}
	pInitCommonControlsEx.Call(uintptr(unsafe.Pointer(&icc)))
	app.queueDir = filepath.Join(os.TempDir(), "PdfMergeUniversal-v313-queue")
	os.MkdirAll(app.queueDir, 0755)
	mutexName := wstr("Local\\PdfMergeUniversal_v313_Mutex")
	mh, _, merr := pCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(mutexName)))
	already := false
	if errno, ok := merr.(syscall.Errno); ok && errno == ERROR_ALREADY_EXISTS {
		already = true
	}
	if already {
		queueSecondary(os.Args[1:])
		if mh != 0 {
			pCloseHandle.Call(mh)
		}
		return
	}
	if mh != 0 {
		defer pCloseHandle.Call(mh)
	}

	app.bgBrush, _, _ = pCreateSolidBrush.Call(rgb(247, 248, 250))
	app.cardBrush, _, _ = pCreateSolidBrush.Call(rgb(255, 255, 255))
	hinst, _, _ := pGetModuleHandleW.Call(0)
	cur, _, _ := pLoadCursorW.Call(0, IDC_ARROW)
	ico, _, _ := pLoadIconW.Call(0, IDI_APPLICATION)
	icoSmall := ico
	if exePath, err := os.Executable(); err == nil {
		iconPath := filepath.Join(filepath.Dir(exePath), "PdfMergeUniversal.ico")
		if h, _, _ := pLoadImageW.Call(0, uintptr(unsafe.Pointer(wstr(iconPath))), IMAGE_ICON, 32, 32, LR_LOADFROMFILE); h != 0 {
			ico = h
		}
		if h, _, _ := pLoadImageW.Call(0, uintptr(unsafe.Pointer(wstr(iconPath))), IMAGE_ICON, 16, 16, LR_LOADFROMFILE); h != 0 {
			icoSmall = h
		}
	}
	cls := wstr("PdfMergeUniversalV313")
	wc := WNDCLASSEX{CbSize: uint32(unsafe.Sizeof(WNDCLASSEX{})), LpfnWndProc: syscall.NewCallback(wndProc), HInstance: hinst, HIcon: ico, HCursor: cur, HbrBackground: app.bgBrush, LpszClassName: cls, HIconSm: icoSmall}
	if r, _, _ := pRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
		return
	}
	hwnd, _, _ := pCreateWindowExW.Call(0, uintptr(unsafe.Pointer(cls)), uintptr(unsafe.Pointer(wstr(appTitle))), WS_OVERLAPPEDWINDOW|WS_VISIBLE, CW_USEDEFAULT, CW_USEDEFAULT, 1160, 840, 0, 0, hinst, 0)
	if hwnd == 0 {
		return
	}
	app.hwnd = hwnd
	pShowWindow.Call(hwnd, SW_SHOW)
	pUpdateWindow.Call(hwnd)
	if len(os.Args) > 1 {
		addPaths(os.Args[1:])
	}
	accel := createAccelerators()
	defer pDestroyAcceleratorTable.Call(accel)
	var m MSG
	for {
		r, _, _ := pGetMessageW.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
		if int32(r) <= 0 {
			break
		}
		if accel != 0 {
			tr, _, _ := pTranslateAcceleratorW.Call(hwnd, accel, uintptr(unsafe.Pointer(&m)))
			if tr != 0 {
				continue
			}
		}
		pTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		pDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
	}
	if app.fontSmall != 0 {
		pDeleteObject.Call(app.fontSmall)
	}
	if app.font != 0 {
		pDeleteObject.Call(app.font)
	}
	if app.fontBold != 0 {
		pDeleteObject.Call(app.fontBold)
	}
	if app.fontTitle != 0 {
		pDeleteObject.Call(app.fontTitle)
	}
	if app.cardBrush != 0 {
		pDeleteObject.Call(app.cardBrush)
	}
	if app.bgBrush != 0 {
		pDeleteObject.Call(app.bgBrush)
	}
}

var _ io.Reader
