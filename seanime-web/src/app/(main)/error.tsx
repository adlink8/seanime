import { LuffyError } from "@/components/shared/luffy-error"
import { Button } from "@/components/ui/button"
import { t } from "@/lib/i18n"
import React from "react"

export default function Error({
    error,
    reset,
}: {
    error: Error & { digest?: string }
    reset: () => void
}) {
    React.useEffect(() => {
        console.error(error)
    }, [error])

    return (
        <div className="flex justify-center">
            <LuffyError
                title={t("common.error.client_side")}
            >
                <p className="max-w-xl text-sm text-[--muted] mb-4">
                    {error.message || t("common.error.unexpected")}
                </p>
                <Button
                    onClick={
                        () => reset()
                    }
                >
                    {t("common.action.retry")}
                </Button>
            </LuffyError>
        </div>
    )
}
