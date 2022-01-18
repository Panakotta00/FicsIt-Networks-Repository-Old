export function setTextRange(el : HTMLElement, start : number, end : number) : void {
    el.focus();
    if (typeof window.getSelection != "undefined" && typeof document.createRange != "undefined") {
        const range = document.createRange();
        range.selectNodeContents(el);
        range.collapse(false);
        range.setStart(el.firstChild, start)
        range.setEnd(el.firstChild, end)
        const sel = window.getSelection();
        sel.removeAllRanges();
        sel.addRange(range);
    }
}

export function placeCaretAtEnd(el : HTMLElement) : void {
    el.focus();
    if (typeof window.getSelection != "undefined" && typeof document.createRange != "undefined") {
        const range = document.createRange();
        range.selectNodeContents(el);
        range.collapse(false);
        const sel = window.getSelection();
        sel.removeAllRanges();
        sel.addRange(range);
    } /*else if (typeof (document.body as HTMLInputElement).createTextRange != "undefined") {
        const textRange = document.body.createTextRange();
        textRange.moveToElementText(el);
        textRange.collapse(false);
        textRange.select();
    }*/
}