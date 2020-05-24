export interface IRouterProps {
    history:  History;
    location: Location;
    match:    Match;
}

export interface History {
    length:   number;
    action:   string;
    location: Location;
    push: (path: string, state?: any)​ => void;
    block: (prompt: any) => void;
    createHref: (location: string) => void;
    go: (n: any) => void;
    goBack: () => void;
    goForward: () => void;
    listen: (listener: any) => void;

    replace: (path: string, state?: any) => void;
}

export interface Location {
    pathname: string;
    search:   string;
    hash:     string;
    key:      string;
    state:    any;
}

export interface Match {
    path:    string;
    url:     string;
    isExact: boolean;
    params:  Params;
}

export interface Params {
}
