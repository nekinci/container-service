

export class AppUtil {


    static GetCurrentApp(): string | null{
        let currentApp: string | null = null;
        if (typeof window !== "undefined"){
            currentApp = localStorage.getItem('currentApp')
        }

        return currentApp;
    }

    static SetCurrentApp(name: string) {
        if (typeof window !== "undefined"){
            localStorage.setItem('currentApp', name)
        }
    }

    static IsThereApp(name: string): boolean {
        let result = false;
        if (typeof  window !== 'undefined'){
            const app = localStorage.getItem('currentApp');
            result = app !== null && app === name;
        }
        return result;
    }

    static IsThereAnyApp(): boolean {
        let result = false;
        if (typeof window !== 'undefined'){
            const app = localStorage.getItem('currentApp');
            result = app !== null;
        }

        return result;
    }

    static DeleteCurrentApp(){
        if (typeof window !== 'undefined'){
            localStorage.removeItem('currentApp');
        }
    }
}