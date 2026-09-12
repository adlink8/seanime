import { SeaImage } from "@/components/shared/sea-image"
import { Button } from "@/components/ui/button/button"
import { cn } from "@/components/ui/core/styling"
import { t } from "@/lib/i18n"
import React from "react"

interface LuffyErrorProps {
    children?: React.ReactNode
    className?: string
    reset?: () => void
    title?: string | null
    showRefreshButton?: boolean
    imageContainerClass?: string
}

export const LuffyError: React.FC<LuffyErrorProps> = (props) => {

    const { children, reset, className, title = t("common.error.title"), showRefreshButton = false, imageContainerClass, ...rest } = props


    return (
        <>
            <div data-luffy-error className={cn("w-full flex flex-col items-center mt-10 space-y-4", className)}>
                {<div
                    data-luffy-error-image-container
                    className={cn("size-[8rem] mx-auto flex-none rounded-[--radius-md] object-cover object-center relative overflow-hidden",
                        imageContainerClass)}
                >
                    <SeaImage
                        data-luffy-error-image
                        src="/luffy-01.png"
                        alt={""}
                        fill
                        priority
                        sizes="10rem"
                        className="object-contain object-top"
                    />
                </div>}
                <div data-luffy-error-content className="text-center space-y-4">
                    {!!title && <h3 data-luffy-error-title>{title}</h3>}
                    <div data-luffy-error-content-children>{children}</div>
                    <div data-luffy-error-content-buttons>
                        {(showRefreshButton && !reset) && (
                            <Button
                                data-luffy-error-content-button-refresh
                                intent="warning-subtle"
                                onClick={() => window.location.reload()}
                            >{t("common.action.retry")}</Button>
                        )}
                        {!!reset && (
                            <Button data-luffy-error-content-button-reset intent="warning-subtle" onClick={reset}>{t("common.action.retry")}</Button>
                        )}
                    </div>
                </div>
            </div>
        </>
    )

}
