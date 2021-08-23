let i = 0;
export enum Roles {
    Scraper = 1 << i++,
    InformationObtainer = 1 << i++,
    Controller = 1 << i++,
    Admin = 1 << i++,
}
