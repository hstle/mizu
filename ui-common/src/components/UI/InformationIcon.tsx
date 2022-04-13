import React, { CSSProperties } from "react";
import styles from "./style/InformationIcon.module.sass"

const DEFUALT_LINK = "https://getmizu.io/docs"

interface LinkProps {
    link?: string,
    className?: string
    title?: string
}

export const Link: React.FC<LinkProps> = ({ link, className, title, children }) => {
    return <React.Fragment>
        <a href={DEFUALT_LINK ? DEFUALT_LINK : link} className={className} title={title} target="_blank">
            {children}
        </a>
    </React.Fragment>
}

export const InformationIcon: React.FC<LinkProps> = ({ link, className }) => {
    return <React.Fragment>
        <Link title="documentation" className={`${styles.linkStyle} ${className}`}>
            <span>Docs</span>
        </Link>
    </React.Fragment>
}