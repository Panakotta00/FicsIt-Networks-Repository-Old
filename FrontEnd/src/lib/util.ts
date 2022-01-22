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

export function getCookie(cookieName : string) : string {
    const cookie = {};
    document.cookie.split(';').forEach(function(el) {
        const [key,value] = el.split('=');
        cookie[key.trim()] = value;
    })
    return cookie[cookieName];
}

export function setCookie(cname : string, cvalue : string, exdays : number) {
    const d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    const expires = "expires="+ d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}
