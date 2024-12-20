export type Message = {
    id: number
    chatUid: string
    type: string
    content: string
    sentAt: string
    changedAt: string
}


export class Api {
    baseUrl: string

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl
    }

    async sendMessage(data: FormData) {
        const res = await fetch(`${this.baseUrl}/api/messages`, {
            method: "POST",
            body: data
        })
        if(!res.ok) {
            throw Error("Request failed")
        }
    }

    listenForMessages(chatUid: string, handler: (e: Message) => void, afterDate?: string): ChatEventListener {
        return new ChatEventListener(this, chatUid, handler, afterDate)
    }

}

export class ChatEventListener {
    private api: Api
    private eventSrc: EventSource
    private afterDate?: string
    private handler: (e: Message) => void 
    private chatUid: string
    constructor(api: Api, chatUid:string, handler: (e: Message) => void, afterDate?: string) {
        this.api = api
        this.afterDate = afterDate
        this.chatUid = chatUid
        this.eventSrc = this.getEventListener()
        this.handler = handler

        this.setupHandlers();
    }

    private setupHandlers() {
        this.eventSrc.onerror = () => {
            this.eventSrc.close()
            this.eventSrc = this.getEventListener()
            this.setupHandlers()
        }
        this.eventSrc.onmessage = (event) => {
            const msg = JSON.parse(event.data) as Message
            if(!this.afterDate || msg.changedAt > this.afterDate) {
                this.afterDate = msg.changedAt
            }
            this.handler(msg)
        }
    }


    close() {
        this.eventSrc.close()
    }


    getEventListener(): EventSource {
        let url = `${this.api.baseUrl}/api/events/messages?chat_uid=${this.chatUid}`
        if(this.afterDate) {
            url += `&after_date=${this.afterDate}`
        }

        const eventSrc = new EventSource(url)
        return eventSrc
    }
}

export let api = new Api("")
