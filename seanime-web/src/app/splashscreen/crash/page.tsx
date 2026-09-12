import { ElectronCrashScreenError } from "@/app/(main)/_electron/electron-crash-screen"
import { LuffyError } from "@/components/shared/luffy-error"
import { LoadingOverlay } from "@/components/ui/loading-spinner"
import { t } from "@/lib/i18n"
import { __isElectronDesktop__ } from "@/types/constants"
import React from "react"

export default function Page() {

    return (
        <LoadingOverlay showSpinner={false}>
            <LuffyError title={t("misc.crash.title")}>
                {__isElectronDesktop__ && <ElectronCrashScreenError />}
            </LuffyError>
        </LoadingOverlay>
    )

}
