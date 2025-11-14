export namespace models {
	
	export class MonitorRule {
	    id: string;
	    name: string;
	    description: string;
	    url: string;
	    method: string;
	    headers?: Record<string, string>;
	    body?: string;
	    interval: number;
	    extractor_type: string;
	    extractor_expr: string;
	    notify_enabled: boolean;
	    enabled: boolean;
	    last_content: string;
	    // Go type: time
	    last_checked: any;
	    status: string;
	    error_message?: string;
	
	    static createFrom(source: any = {}) {
	        return new MonitorRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.url = source["url"];
	        this.method = source["method"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	        this.interval = source["interval"];
	        this.extractor_type = source["extractor_type"];
	        this.extractor_expr = source["extractor_expr"];
	        this.notify_enabled = source["notify_enabled"];
	        this.enabled = source["enabled"];
	        this.last_content = source["last_content"];
	        this.last_checked = this.convertValues(source["last_checked"], null);
	        this.status = source["status"];
	        this.error_message = source["error_message"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

