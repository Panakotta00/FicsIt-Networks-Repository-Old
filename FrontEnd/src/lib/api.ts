export class Package_Base {
    id = -1
    name = ""
    displayName = ""
    description = ""
    sourceLink = ""
    creatorId = -1
}

export class Package extends Package_Base {
    verified = false
}

export class Tag_Base {
    id : bigint
    name : string
    description : string

    constructor(v : {id?:bigint, name?:string, description?:string}) {
        this.id = (v && v.id) || -1 as unknown as bigint;
        this.name = (v && v.name) || "New Tag";
        this.description = (v && v.description) || "New Tag Description...";
    }
}

export class Tag extends Tag_Base {
    verified = false;
}

export class User_Base {
    id = -1
    name = ""
    bio = ""
    admin = false
}

export class User extends User_Base {
    verified = false
}

export type PackageRef = bigint|string;

export function apiURL(endpoint : string) : string {
    return new URL(endpoint, import.meta.env.VITE_API_BASE_URL as string).toString()
}

export const listPackages = async (page : number, count : number) : Promise<Package_Base[]> => {
    const res = await fetch(apiURL(`package?page=${page}&count=${count}`))
    if (!res.ok) throw new Error('Bad response')
    const items = await res.json()
    return items as Package[]
}

export const getPackage = async (packageRef : PackageRef) : Promise<Package_Base> => {
    const res = await fetch(apiURL(`package/${packageRef}`))
    if (!res.ok) throw new Error('Bad response')
    const pack = await res.json()
    return pack as Package
}

export const getPackageTags = async (packageId : bigint) : Promise<Tag_Base[]> => {
    const res = await fetch(apiURL(`package/${packageId}/tags`))
    if (!res.ok) throw new Error('Bad response')
    const tag = await res.json()
    return tag as Tag[]
}

export const getUser = async(userId : bigint) : Promise<User> => {
    const res = await fetch(apiURL(`user/${userId}`))
    if (!res.ok) throw new Error('Bad response')
    const user = await res.json()
    return user as User
}

export const getTags = async() : Promise<Tag_Base[]> => {
    const res = await fetch(apiURL(`tag`))
    if (!res.ok) throw new Error('Bad response')
    const tag = await res.json()
    return tag as Tag[]
}