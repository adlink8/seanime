import { Button, ButtonProps } from "@/components/ui/button"
import { Modal } from "@/components/ui/modal"
import { useDisclosure, UseDisclosureReturn } from "@/hooks/use-disclosure"
import { t } from "@/lib/i18n"
import React from "react"

type ConfirmationDialogHookProps = {
    title: string,
    description?: React.ReactNode,
    actionText?: string,
    cancelText?: string,
    actionIntent?: ButtonProps["intent"]
    onConfirm: () => void
    onCancel?: () => void
}

export function useConfirmationDialog(props: ConfirmationDialogHookProps) {
    const api = useDisclosure(false)
    return {
        ...api,
        ...props,
    }
}

export const ConfirmationDialog: React.FC<ConfirmationDialogHookProps & UseDisclosureReturn> = (props) => {

    const {
        isOpen,
        close,
        onConfirm,
        onCancel,
        title,
        description = t("common.dialog.confirm_description"),
        actionText = t("common.action.confirm"),
        cancelText = t("common.action.cancel"),
        actionIntent = "alert-subtle",
    } = props

    return (
        <>
            <Modal
                title={title}
                titleClass="text-center"
                open={isOpen}
                onOpenChange={close}
            >
                <div className="space-y-4">
                    <div className="text-center">{description}</div>
                    <div className="flex gap-2 justify-center items-center">
                        <Button
                            intent={actionIntent}
                            onClick={() => {
                                onConfirm()
                                close()
                            }}
                        >
                            {actionText}
                        </Button>
                        <Button
                            intent="white"
                            onClick={() => {
                                onCancel?.()
                                close()
                            }}
                        >
                            {cancelText}
                        </Button>
                    </div>
                </div>
            </Modal>
        </>
    )

}
