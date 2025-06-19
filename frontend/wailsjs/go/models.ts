export namespace config {
	
	export class Config {
	    WorkspaceRoot: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WorkspaceRoot = source["WorkspaceRoot"];
	    }
	}

}

