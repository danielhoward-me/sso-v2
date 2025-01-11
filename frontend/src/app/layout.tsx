import "@radix-ui/themes/styles.css";

import { Theme } from "@radix-ui/themes";
import { headers } from "next/headers";

const ACCOUNT_HOSTNAME = process.env.ACCOUNT_HOSTNAME;
const SSO_HOSTNAME = process.env.SSO_HOSTNAME;

export default async function RootLayout({ children, account, sso }: {
    children: React.ReactNode,
    account?: React.ReactNode,
    sso?: React.ReactNode,
}) {
    const headersList = await headers();
    // Header set by nginx
    const host = headersList.get('Host');

    let page = children;

    switch (host) {
        case ACCOUNT_HOSTNAME: {
            if (account !== undefined) page = account;
            break;
        }
        case SSO_HOSTNAME: {
            if (sso !== undefined) page = sso;
            break;
        }
        default: {
            throw new Error(`'${host}' is not a valid host`);
        }
    }

    return (
        <html lang="en">
            <body>
                <Theme>
                    {page}
                </Theme>
            </body>
        </html>
    )
}
