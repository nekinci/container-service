
class Event {
    list: Map<string, any[]> = new Map<string, any[]>();

    on(eventType: string, cb: any) {
        if (!this.list.has(eventType)) {
            this.list.set(eventType, []);
        }

        // @ts-ignore
        this.list.get(eventType).push(cb);
    }

    emit(eventType: string, ...args: any) {
        if (this.list.has(eventType)){
            // @ts-ignore
            this?.list?.get(eventType).forEach((cb) => {
                cb(...args);
            })
        } else {
            console.error("event type not allowed!")
        }
    }
}

const event = new Event();

export default event;