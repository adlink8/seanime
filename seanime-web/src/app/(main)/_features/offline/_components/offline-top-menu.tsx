import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { NavigationMenu, NavigationMenuProps } from "@/components/ui/navigation-menu"
import { t } from "@/lib/i18n"
import { usePathname } from "@/lib/navigation"
import React, { useMemo } from "react"

interface OfflineTopMenuProps {
    children?: React.ReactNode
}

export const OfflineTopMenu: React.FC<OfflineTopMenuProps> = (props) => {

    const { children, ...rest } = props

    const serverStatus = useServerStatus()

    const pathname = usePathname()

    const navigationItems = useMemo<NavigationMenuProps["items"]>(() => {

        return [
            {
                href: "/offline",
                // icon: IoLibrary,
                isCurrent: pathname === "/offline",
                name: t("navigation.offline.anime_library"),
            },
            ...[serverStatus?.settings?.library?.enableManga && {
                href: "/offline/manga",
                icon: null,
                isCurrent: pathname.includes("/offline/manga"),
                name: t("navigation.item.manga"),
            }].filter(Boolean) as NavigationMenuProps["items"],
        ].filter(Boolean)
    }, [pathname, serverStatus?.settings?.library?.enableManga])

    return (
        <NavigationMenu
            className="p-0 hidden lg:inline-block"
            itemClass="text-xl"
            items={navigationItems}
        />
    )

}
