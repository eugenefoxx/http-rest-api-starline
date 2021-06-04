let resizeObserver = null;
console.log('I AM MODAL');

const CLASS_LIST = {
    MODAL: 'modal',
    MODAL_ACTIVE: 'modal--active',
    MODAL_HAS_SCROLL: 'modal--has-sc',
    MODAL_DIALOG_BODY: 'modal__dialog-body',
    TRIGGER_OPEN: 'js-modal-open',
    TRIGGER_CLOSE: 'js-modal-close'
};

const showScroll = (event) => {
    if (event.propertyName === 'transform') {
        document.body.style.paddingRight = ``;
        document.body.style.overflow = 'visible'; 
        
        event.target.closest(`.${CLASS_LIST.MODAL}`).removeEventListener('transitionend', showScroll);
    }
}

//debugger;

document.addEventListener('click', (event) => {
    // open
   // debugger;

    if (event.target.closest(`.${CLASS_LIST.TRIGGER_OPEN}`)) {
        console.log('open');
        event.preventDefault();

        const target = event.target.closest(`.${CLASS_LIST.TRIGGER_OPEN}`);
        const modalId = target.getAttribute('href').replace('#', '');
        const modal = document.getElementById(modalId);

    //    document.body.style.paddingRight = `${getScrollbarWidth()}px`;
        document.body.style.overflow = 'hidden';

        modal.classList.add(CLASS_LIST.MODAL_ACTIVE);

        bindResizeObserver(modal);
    }
    // close
    if (
        event.target.closest(`.${CLASS_LIST.TRIGGER_CLOSE}`) ||
        event.target.classList.contains(CLASS_LIST.MODAL_ACTIVE)
    ) {
        console.log('close');
        event.preventDefault();

        const modal = event.target.closest(`.${CLASS_LIST.MODAL}`);

        modal.classList.remove(CLASS_LIST.MODAL_ACTIVE);

        unbindResizeObserver(modal);

        modal.addEventListener('transitionend', showScroll);
    }
});

const getScrollbarWidth = () => {
    const item = document.createElement('div');

    item.style.position = 'absolute';
    item.style.top = '-9999px';
    item.style.width = '50px';
    item.style.height = '50px';
    item.style.overflow = 'scroll';
    item.style.visibility = 'hidden';

    document.body.appendChild(item);
    const scrollBarWidth = item.offsetWidth - item.clientWidth;
    document.body.removeChild(item);

    return scrollBarWidth;

};

const bindResizeObserver = (modal) => {
    const content = modal.querySelector(`.${CLASS_LIST.MODAL_DIALOG_BODY}`);

    const toggleShadows = () => {
        modal.classList.toggle(
            CLASS_LIST.MODAL_HAS_SCROLL,
            content.scrollHeight > content.clientHeight
        );
    };

    resizeObserver = new ResizeObserver(toggleShadows);

    resizeObserver.observe(content);
};

const unbindResizeObserver = (modal) => {
    const content = modal.querySelector(`.${CLASS_LIST.MODAL_DIALOG_BODY}`);

    resizeObserver.unobserve(content);
    resizeObserver = null;
};
/*
function validateQrCode(code) {
    const regexp = /\bP\d{7}LK\d{9}R\d{10}Q\d{5}D\d{8}\b/;
    const regexp2 = /\bP\d{7}L\d{10}R\d{10}Q\d{5}D\d{8}\b/;

    return (
        (regexp.test(code) || regexp2.test(code))
        && document.querySelector('td[data-material-id="'+code+'"]') !== null
    );
}

function onQrCodeChange(e) {
    const qrCode = e.target.value;

    if (event.keyCode === 13 && validateQrCode(qrCode)) {
        // 1. 
        // 2.
    }
    
}

document.addEventListener("DOMContentLoaded", function() {
    document.getElementById('qr-code-value').onkeyup = onQrCodeChange;
});
*/
