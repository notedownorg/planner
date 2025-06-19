export namespace config {
	
	export class WeeklyViewComponents {
	    HabitTracker: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WeeklyViewComponents(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.HabitTracker = source["HabitTracker"];
	    }
	}
	export class WeeklyViewConfig {
	    EnabledComponents: WeeklyViewComponents;
	
	    static createFrom(source: any = {}) {
	        return new WeeklyViewConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EnabledComponents = this.convertValues(source["EnabledComponents"], WeeklyViewComponents);
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
	export class PeriodicNotes {
	    WeeklySubdir: string;
	    WeeklyNameFormat: string;
	
	    static createFrom(source: any = {}) {
	        return new PeriodicNotes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WeeklySubdir = source["WeeklySubdir"];
	        this.WeeklyNameFormat = source["WeeklyNameFormat"];
	    }
	}
	export class Config {
	    WorkspaceRoot: string;
	    PeriodicNotes: PeriodicNotes;
	    WeeklyView: WeeklyViewConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.WorkspaceRoot = source["WorkspaceRoot"];
	        this.PeriodicNotes = this.convertValues(source["PeriodicNotes"], PeriodicNotes);
	        this.WeeklyView = this.convertValues(source["WeeklyView"], WeeklyViewConfig);
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

export namespace habits {
	
	export class Habit {
	    name: string;
	    completed: boolean;
	    order: number;
	
	    static createFrom(source: any = {}) {
	        return new Habit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.completed = source["completed"];
	        this.order = source["order"];
	    }
	}
	export class WeeklyHabits {
	    year: number;
	    week_number: number;
	    habits: Record<string, Habit>;
	    day_status: Record<string, boolean>;
	
	    static createFrom(source: any = {}) {
	        return new WeeklyHabits(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.year = source["year"];
	        this.week_number = source["week_number"];
	        this.habits = this.convertValues(source["habits"], Habit, true);
	        this.day_status = source["day_status"];
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

