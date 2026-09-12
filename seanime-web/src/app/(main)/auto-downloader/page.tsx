import { CustomLibraryBanner } from "@/app/(main)/_features/anime-library/_containers/custom-library-banner"
import { AutoDownloaderPage } from "@/app/(main)/auto-downloader/_containers/autodownloader-page"
import { PageWrapper } from "@/components/shared/page-wrapper"
import { t } from "@/lib/i18n"
import React from "react"


export default function Page() {

    return (
        <>
            <CustomLibraryBanner discrete />
            <PageWrapper className="p-4 sm:p-8 space-y-4">
                <div className="flex justify-between items-center w-full relative">
                    <div>
                        <h2>{t("settings.nav.auto_downloader_menu")}</h2>
                        <p className="text-[--muted]">
                            {t("autodownloader.page.subtitle")}
                        </p>
                    </div>
                </div>
                <AutoDownloaderPage />
            </PageWrapper>
        </>
    )

}
