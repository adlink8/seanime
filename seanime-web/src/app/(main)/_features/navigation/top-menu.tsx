import { useMissingEpisodeCount } from "@/app/(main)/_hooks/missing-episodes-loader"
import { useServerStatus } from "@/app/(main)/_hooks/use-server-status"
import { NavigationMenu, NavigationMenuProps } from "@/components/ui/navigation-menu"
import { usePathname } from "@/lib/navigation"
import React, { useMemo } from "react"

interface TopMenuProps {
    children?: React.ReactNode
}

export const TopMenu: React.FC<TopMenuProps> = (props) => {

    const { children, ...rest } = props

    const serverStatus = useServerStatus()

    const pathname = usePathname()

    const missingEpisodeCount = useMissingEpisodeCount()

    const navigationItems = useMemo<NavigationMenuProps["items"]>(() => {

        return [
            {
                href: "/",
                // icon: IoLibrary,
                isCurrent: pathname === "/",
                name: "首页",
            },
            {
                href: "/schedule",
                icon: null,
                isCurrent: pathname.startsWith("/schedule"),
                name: "放送日历",
                // addon: missingEpisodeCount > 0 ? <Badge
                //     className="absolute -top-1 right-2 h-2 w-2 p-0 z-[5]" size="sm"
                //     intent="alert-solid"
                // /> : undefined,
            },
            ...[serverStatus?.settings?.library?.enableManga && {
                href: "/manga",
                icon: null,
                isCurrent: pathname.startsWith("/manga"),
                name: "漫画",
            }].filter(Boolean) as NavigationMenuProps["items"],
            {
                href: "/lists",
                icon: null,
                isCurrent: pathname.startsWith("/lists"),
                name: "我的片单",
            },
            {
                href: "/discover",
                icon: null,
                isCurrent: pathname.startsWith("/discover") || pathname.startsWith("/search"),
                name: "探索发现",
            },
        ].filter(Boolean)
    }, [pathname, missingEpisodeCount, serverStatus?.settings?.library?.enableManga])

    return (
        <NavigationMenu
            className="p-0 hidden lg:inline-block"
            items={navigationItems}
            desktopListClass="space-x-0"
            data-top-menu
        />
    )

}
