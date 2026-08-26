
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/actions" | "/actions/new" | "/actions/[id]" | "/pages";
		RouteParams(): {
			"/actions/[id]": { id: string }
		};
		LayoutParams(): {
			"/": { id?: string };
			"/actions": { id?: string };
			"/actions/new": Record<string, never>;
			"/actions/[id]": { id: string };
			"/pages": Record<string, never>
		};
		Pathname(): "/" | "/actions" | "/actions/" | "/actions/new" | "/actions/new/" | `/actions/${string}` & {} | `/actions/${string}/` & {} | "/pages" | "/pages/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): "/favicon.svg" | "/robots.txt" | string & {};
	}
}