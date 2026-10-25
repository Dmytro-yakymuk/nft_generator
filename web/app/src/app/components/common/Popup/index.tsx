import { Dispatch, SetStateAction, useEffect, useRef } from 'react';

import './index.scss';

/** Popup layer is custom layer with props for popups. */
export const Popup: React.FC<{ setIsPopupVisible: Dispatch<SetStateAction<boolean>> }> = ({
    children,
    setIsPopupVisible,
}) => {
    const popupRef = useRef<HTMLDivElement | null>(null);

    /** Changes popup visability. */
    const changePopupVisability = (e: any) => {
        if (!popupRef.current) {
            return;
        };

        setIsPopupVisible(popupRef.current.contains(e.target));
    };

    useEffect(() => {
        document.body.addEventListener(
            'click',
            changePopupVisability
        );

        return () => {
            document.body.removeEventListener(
                'click',
                changePopupVisability
            );
        };
    });

    return <section className="popup">
        <div
            className="popup__wrapper"
            ref={popupRef}
        >
            {children}
        </div>
    </section>;
};
