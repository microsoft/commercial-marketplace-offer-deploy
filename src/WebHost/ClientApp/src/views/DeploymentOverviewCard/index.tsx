import React from "react"
import { Image } from "react-bootstrap";
import "./style.scss";
import classNames from 'classnames';

export const OverViewCard = (props: OverViewCardProps) => {
    return (
        <div onClick={() => props.onClick(props.value)} className={
            classNames({
                [props.className]: true,
                'overview-card': true,
                'active': props.value == props.selectedFilter
            })
        }>
            <div className="overview-card-container">
                <div className="overview-item">
                    <Image src={props.image} style={{ backgroundColor: props.color }} ></Image>
                </div>
                <div className="overview-item-conent ">
                    <span className="overview-card-title">{props.title}</span>
                    <span className="overview-card-sub-title">{props.subTitle}</span>
                </div>
            </div>
        </div>
    )
}

interface OverViewCardProps {
    className: any
    image: any
    title?: string
    subTitle: string
    value: string
    color: string
    selectedFilter: string
    onClick(filter: string): any
}