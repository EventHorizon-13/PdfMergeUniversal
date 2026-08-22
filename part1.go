//go:build windows

package main

import (
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func wstr(s string) *uint16 { p, _ := syscall.UTF16PtrFromString(s); return p }

func utf16Multi(s string) []uint16 {
	if !strings.HasSuffix(s, "\x00\x00") { s += "\x00\x00" }
	return utf16.Encode([]rune(s))
}
func loword(v uintptr) uint16 { return uint16(v & 0xffff) }
func hiword(v uintptr) uint16 { return uint16((v >> 16) & 0xffff) }
func boolPtr(b bool) uintptr { if b { return 1 }; return 0 }

func msgbox(title, text string, flags uintptr) int {
	r, _, _ := pMessageBoxW.Call(app.hwnd, uintptr(unsafe.Pointer(wstr(text))), uintptr(unsafe.Pointer(wstr(title))), flags)
	return int(r)
}
func setText(hwnd uintptr, s string) { pSetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(wstr(s)))) }
func getText(hwnd uintptr) string {
	n, _, _ := pGetWindowTextLengthW.Call(hwnd); if n == 0 { return "" }
	buf := make([]uint16, n+1); pGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), n+1); return syscall.UTF16ToString(buf)
}
func send(hwnd uintptr, msg uint32, wp, lp uintptr) uintptr { r, _, _ := pSendMessageW.Call(hwnd, uintptr(msg), wp, lp); return r }
func move(hwnd uintptr, x, y, w, h int32) { if hwnd != 0 { pMoveWindow.Call(hwnd, uintptr(x), uintptr(y), uintptr(w), uintptr(h), 1) } }
func enable(hwnd uintptr, b bool) { pEnableWindow.Call(hwnd, boolPtr(b)) }

func createControl(exStyle uint32, class, text string, style uint32, id int32) uintptr {
	hinst, _, _ := pGetModuleHandleW.Call(0)
	h, _, _ := pCreateWindowExW.Call(uintptr(exStyle), uintptr(unsafe.Pointer(wstr(class))), uintptr(unsafe.Pointer(wstr(text))), uintptr(style|WS_CHILD|WS_VISIBLE), 0, 0, 0, 0, app.hwnd, uintptr(id), hinst, 0)
	if h != 0 && app.font != 0 { send(h, WM_SETFONT, app.font, 1) }
	if h != 0 { pSetWindowTheme.Call(h, uintptr(unsafe.Pointer(wstr("Explorer"))), 0) }
	return h
}
func setFont(hwnd, font uintptr) { send(hwnd, WM_SETFONT, font, 1) }
func createFonts() {
	var lf LOGFONT
	ok, _, _ := pSystemParametersInfoW.Call(SPI_GETICONTITLELOGFONT, uintptr(unsafe.Sizeof(lf)), uintptr(unsafe.Pointer(&lf)), 0)
	if ok == 0 { face := syscall.StringToUTF16("Segoe UI"); copy(lf.FaceName[:], face); lf.Height, lf.Weight, lf.CharSet, lf.Quality = -12, 400, 1, 5 }
	app.font, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&lf)))
	small := lf; app.fontSmall, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&small)))
	bold := lf; bold.Weight = 600; app.fontBold, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&bold)))
	title := lf; title.Height = title.Height * 3 / 2; title.Weight = 650; app.fontTitle, _, _ = pCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&title)))
}
func rgb(r, g, b byte) uintptr { return uintptr(r) | uintptr(g)<<8 | uintptr(b)<<16 }

func drawButton(di *DRAWITEMSTRUCT) {
	if di == nil { return }
	if di.HwndItem == app.outputCard {
		br, _, _ := pCreateSolidBrush.Call(rgb(255,255,255)); pen,_,_ := pCreatePen.Call(PS_SOLID,1,rgb(234,236,240)); oldBr,_,_ := pSelectObject.Call(di.HDC,br); oldPen,_,_ := pSelectObject.Call(di.HDC,pen)
		pRoundRect.Call(di.HDC,uintptr(di.RcItem.Left),uintptr(di.RcItem.Top),uintptr(di.RcItem.Right),uintptr(di.RcItem.Bottom),14,14); pSelectObject.Call(di.HDC,oldBr); pSelectObject.Call(di.HDC,oldPen); pDeleteObject.Call(br); pDeleteObject.Call(pen); return
	}
	primary := int(di.CtlID)==ID_MERGE; disabled := (di.ItemState&ODS_DISABLED)!=0; pressed := (di.ItemState&ODS_SELECTED)!=0
	bg:=rgb(255,255,255); fg:=rgb(52,64,84); border:=rgb(208,213,221)
	if primary { bg=rgb(37,99,235); fg=rgb(255,255,255); border=rgb(37,99,235); if pressed { bg=rgb(29,78,216); border=bg } } else if pressed { bg=rgb(242,244,247) }
	if disabled { bg=rgb(242,244,247); fg=rgb(152,162,179); border=rgb(228,231,236) }
	br,_,_:=pCreateSolidBrush.Call(bg); pen,_,_:=pCreatePen.Call(PS_SOLID,1,border); oldBr,_,_:=pSelectObject.Call(di.HDC,br); oldPen,_,_:=pSelectObject.Call(di.HDC,pen)
	pRoundRect.Call(di.HDC,uintptr(di.RcItem.Left),uintptr(di.RcItem.Top),uintptr(di.RcItem.Right),uintptr(di.RcItem.Bottom),10,10); pSelectObject.Call(di.HDC,oldBr); pSelectObject.Call(di.HDC,oldPen); pDeleteObject.Call(br); pDeleteObject.Call(pen)
	pSetBkMode.Call(di.HDC,TRANSPARENT); pSetTextColor.Call(di.HDC,fg); font:=app.font; if primary { font=app.fontBold }; oldFont,_,_:=pSelectObject.Call(di.HDC,font); text:=getText(di.HwndItem); rr:=di.RcItem
	pDrawTextW.Call(di.HDC,uintptr(unsafe.Pointer(wstr(text))),^uintptr(0),uintptr(unsafe.Pointer(&rr)),DT_CENTER|DT_VCENTER|DT_SINGLELINE); pSelectObject.Call(di.HDC,oldFont)
}
func staticColor(hdc, child uintptr) uintptr {
	pSetBkMode.Call(hdc, TRANSPARENT)
	if child==app.countLabel { pSetTextColor.Call(hdc,rgb(37,99,235)) } else if child==app.subLabel || child==app.statusLabel || child==app.qpdfLabel || child==app.finalLabel { pSetTextColor.Call(hdc,rgb(102,112,133)) } else { pSetTextColor.Call(hdc,rgb(16,24,40)) }
	return app.bgBrush
}
func initListColumns() {
	cols:=[]struct{text string; width int32; fmt int32}{{"序号",55,LVCFMT_LEFT},{"文件名",425,LVCFMT_LEFT},{"页数",65,LVCFMT_RIGHT},{"所在文件夹",400,LVCFMT_LEFT}}
	for i,c:=range cols { txt:=syscall.StringToUTF16(c.text); col:=LVCOLUMN{Mask:LVCF_TEXT|LVCF_WIDTH|LVCF_FMT,Fmt:c.fmt,Cx:c.width,PszText:&txt[0]}; send(app.list,LVM_INSERTCOLUMNW,uintptr(i),uintptr(unsafe.Pointer(&col))) }
	send(app.list,LVM_SETEXTENDEDLISTVIEWSTYLE,0,LVS_EX_FULLROWSELECT|LVS_EX_DOUBLEBUFFER)
}
func insertListRow(index int, item PdfItem) {
	vals:=[]string{strconv.Itoa(index+1),filepath.Base(item.Path),"?",filepath.Dir(item.Path)}; if item.Known { vals[2]=strconv.Itoa(item.Pages) }
	for sub,s:=range vals { txt:=syscall.StringToUTF16(s); li:=LVITEM{Mask:LVIF_TEXT,IItem:int32(index),ISubItem:int32(sub),PszText:&txt[0]}; if sub==0 { send(app.list,LVM_INSERTITEMW,0,uintptr(unsafe.Pointer(&li))) } else { send(app.list,LVM_SETITEMW,0,uintptr(unsafe.Pointer(&li))) } }
}
