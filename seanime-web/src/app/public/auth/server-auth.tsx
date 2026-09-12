import { serverAuthTokenAtom } from "@/app/(main)/_atoms/server-status.atoms"
import { defineSchema, Field, Form } from "@/components/ui/form"
import { Modal } from "@/components/ui/modal"
import { t } from "@/lib/i18n"
import { useAtom } from "jotai"
import { sha256 } from "js-sha256"
import React, { useState } from "react"

export function ServerAuth() {

    const [, setAuthToken] = useAtom(serverAuthTokenAtom)
    const [loading, setLoading] = useState(false)

    return (<>
        <Modal
            title={t("misc.auth.title")}
            description={t("misc.auth.description")}
            open={true}
            onOpenChange={(v) => {}}
            overlayClass="bg-opacity-100 bg-gray-900"
            contentClass="border focus:outline-none focus-visible:outline-none outline-none"
            hideCloseButton
        >
            <Form
                schema={defineSchema(({ z }) => z.object({
                    password: z.string().min(1, t("misc.auth.password_required")),
                }))}
                onSubmit={async data => {
                    setLoading(true)
                    const hash = sha256(data.password)
                    setAuthToken(hash)
                    React.startTransition(() => {
                        window.location.href = "/"
                        setLoading(false)
                    })
                }}
            >
                <Field.Text
                    type="password"
                    name="password"
                    label={t("misc.auth.enter_password")}
                    fieldClass=""
                />
                <Field.Submit showLoadingOverlayOnSuccess loading={loading}>{t("media.action.continue")}</Field.Submit>
            </Form>
        </Modal>
    </>)
}
