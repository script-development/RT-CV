let i = 0;
export enum Roles {
    Scraper = 1 << i++,
    InformationObtainer = 1 << i++,
    Controller = 1 << i++,
    Admin = 1 << i++,
}

export const allRoles: Array<Roles> = [
    Roles.Scraper,
    Roles.InformationObtainer,
    Roles.Controller,
    Roles.Admin,
]

interface RoleInfo {
    title: string
    description: string
}

export function roleInfo(role: Roles): undefined | RoleInfo {
    switch (role) {
        case Roles.Scraper:
            return {
                title: 'Scraper',
                description: 'Can insert scraped data'
            }
        case Roles.InformationObtainer:
            return {
                title: 'Information obtainer',
                description: 'Can obtain information the server has'
            }
        case Roles.Controller:
            return {
                title: 'Controller',
                description: 'Can control the server'
            }
        case Roles.Admin:
            return {
                title: 'Admin',
                description: 'Currently unused role',
            }
    }
}
