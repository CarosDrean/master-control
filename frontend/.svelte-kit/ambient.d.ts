
// this file is generated — do not edit it


/// <reference types="@sveltejs/kit" />

/**
 * Environment variables [loaded by Vite](https://vitejs.dev/guide/env-and-mode.html#env-files) from `.env` files and `process.env`. Like [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), this module cannot be imported into client-side code. This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured).
 * 
 * _Unlike_ [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), the values exported from this module are statically injected into your bundle at build time, enabling optimisations like dead code elimination.
 * 
 * ```ts
 * import { API_KEY } from '$env/static/private';
 * ```
 * 
 * Note that all environment variables referenced in your code should be declared (for example in an `.env` file), even if they don't have a value until the app is deployed:
 * 
 * ```
 * MY_FEATURE_FLAG=""
 * ```
 * 
 * You can override `.env` values from the command line like so:
 * 
 * ```sh
 * MY_FEATURE_FLAG="enabled" npm run dev
 * ```
 */
declare module '$env/static/private' {
	export const WARP_CLI_AGENT_PROTOCOL_VERSION: string;
	export const ZLE_RPROMPT_INDENT: string;
	export const TERM_PROGRAM: string;
	export const NODE: string;
	export const POSH_SHELL: string;
	export const OPENCODE_PID: string;
	export const INIT_CWD: string;
	export const WARP_HONOR_PS1: string;
	export const SHELL: string;
	export const TERM: string;
	export const TMPDIR: string;
	export const WARP_PROMPT_NODE_VERSION_ENABLED: string;
	export const WARP_TERMINAL_SESSION_UUID: string;
	export const TERM_PROGRAM_VERSION: string;
	export const CONDA_PROMPT_MODIFIER: string;
	export const ZSH: string;
	export const POWERLINE_COMMAND: string;
	export const USER: string;
	export const POSH_SESSION_ID: string;
	export const LS_COLORS: string;
	export const COMMAND_MODE: string;
	export const PNPM_SCRIPT_SRC_DIR: string;
	export const SSH_AUTH_SOCK: string;
	export const WARP_IS_LOCAL_SHELL_SESSION: string;
	export const __CF_USER_TEXT_ENCODING: string;
	export const PYENV_VIRTUALENV_DISABLE_PROMPT: string;
	export const WARP_SSH_REUSE_CONTROL_MASTER: string;
	export const npm_execpath: string;
	export const VIRTUAL_ENV_DISABLE_PROMPT: string;
	export const PAGER: string;
	export const WARP_USE_SSH_WRAPPER: string;
	export const LSCOLORS: string;
	export const PATH: string;
	export const LaunchInstanceID: string;
	export const npm_package_json: string;
	export const __CFBundleIdentifier: string;
	export const PWD: string;
	export const npm_command: string;
	export const POSH_CONFIG: string;
	export const npm_lifecycle_event: string;
	export const POSH_SHELL_VERSION: string;
	export const LANG: string;
	export const npm_package_name: string;
	export const WARP_FOCUS_URL: string;
	export const NODE_PATH: string;
	export const XPC_FLAGS: string;
	export const npm_config_node_gyp: string;
	export const XPC_SERVICE_NAME: string;
	export const npm_package_version: string;
	export const pnpm_config_verify_deps_before_run: string;
	export const HOME: string;
	export const SHLVL: string;
	export const OSTYPE: string;
	export const LESS: string;
	export const LOGNAME: string;
	export const npm_lifecycle_script: string;
	export const LC_CTYPE: string;
	export const SSH_SOCKET_DIR: string;
	export const GOPATH: string;
	export const npm_config_user_agent: string;
	export const WARP_CLIENT_VERSION: string;
	export const OSLogRateLimit: string;
	export const AGENT: string;
	export const CONDA_CHANGEPS1: string;
	export const SECURITYSESSIONID: string;
	export const COLORTERM: string;
	export const OPENCODE: string;
	export const npm_node_execpath: string;
	export const NODE_ENV: string;
}

/**
 * Similar to [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private), except that it only includes environment variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`), and can therefore safely be exposed to client-side code.
 * 
 * Values are replaced statically at build time.
 * 
 * ```ts
 * import { PUBLIC_BASE_URL } from '$env/static/public';
 * ```
 */
declare module '$env/static/public' {
	
}

/**
 * This module provides access to runtime environment variables, as defined by the platform you're running on. For example if you're using [`adapter-node`](https://github.com/sveltejs/kit/tree/main/packages/adapter-node) (or running [`vite preview`](https://svelte.dev/docs/kit/cli)), this is equivalent to `process.env`. This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured).
 * 
 * This module cannot be imported into client-side code.
 * 
 * ```ts
 * import { env } from '$env/dynamic/private';
 * console.log(env.DEPLOYMENT_SPECIFIC_VARIABLE);
 * ```
 * 
 * > [!NOTE] In `dev`, `$env/dynamic` always includes environment variables from `.env`. In `prod`, this behavior will depend on your adapter.
 */
declare module '$env/dynamic/private' {
	export const env: {
		WARP_CLI_AGENT_PROTOCOL_VERSION: string;
		ZLE_RPROMPT_INDENT: string;
		TERM_PROGRAM: string;
		NODE: string;
		POSH_SHELL: string;
		OPENCODE_PID: string;
		INIT_CWD: string;
		WARP_HONOR_PS1: string;
		SHELL: string;
		TERM: string;
		TMPDIR: string;
		WARP_PROMPT_NODE_VERSION_ENABLED: string;
		WARP_TERMINAL_SESSION_UUID: string;
		TERM_PROGRAM_VERSION: string;
		CONDA_PROMPT_MODIFIER: string;
		ZSH: string;
		POWERLINE_COMMAND: string;
		USER: string;
		POSH_SESSION_ID: string;
		LS_COLORS: string;
		COMMAND_MODE: string;
		PNPM_SCRIPT_SRC_DIR: string;
		SSH_AUTH_SOCK: string;
		WARP_IS_LOCAL_SHELL_SESSION: string;
		__CF_USER_TEXT_ENCODING: string;
		PYENV_VIRTUALENV_DISABLE_PROMPT: string;
		WARP_SSH_REUSE_CONTROL_MASTER: string;
		npm_execpath: string;
		VIRTUAL_ENV_DISABLE_PROMPT: string;
		PAGER: string;
		WARP_USE_SSH_WRAPPER: string;
		LSCOLORS: string;
		PATH: string;
		LaunchInstanceID: string;
		npm_package_json: string;
		__CFBundleIdentifier: string;
		PWD: string;
		npm_command: string;
		POSH_CONFIG: string;
		npm_lifecycle_event: string;
		POSH_SHELL_VERSION: string;
		LANG: string;
		npm_package_name: string;
		WARP_FOCUS_URL: string;
		NODE_PATH: string;
		XPC_FLAGS: string;
		npm_config_node_gyp: string;
		XPC_SERVICE_NAME: string;
		npm_package_version: string;
		pnpm_config_verify_deps_before_run: string;
		HOME: string;
		SHLVL: string;
		OSTYPE: string;
		LESS: string;
		LOGNAME: string;
		npm_lifecycle_script: string;
		LC_CTYPE: string;
		SSH_SOCKET_DIR: string;
		GOPATH: string;
		npm_config_user_agent: string;
		WARP_CLIENT_VERSION: string;
		OSLogRateLimit: string;
		AGENT: string;
		CONDA_CHANGEPS1: string;
		SECURITYSESSIONID: string;
		COLORTERM: string;
		OPENCODE: string;
		npm_node_execpath: string;
		NODE_ENV: string;
		[key: `PUBLIC_${string}`]: undefined;
		[key: `${string}`]: string | undefined;
	}
}

/**
 * Similar to [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private), but only includes variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`), and can therefore safely be exposed to client-side code.
 * 
 * Note that public dynamic environment variables must all be sent from the server to the client, causing larger network requests — when possible, use `$env/static/public` instead.
 * 
 * ```ts
 * import { env } from '$env/dynamic/public';
 * console.log(env.PUBLIC_DEPLOYMENT_SPECIFIC_VARIABLE);
 * ```
 */
declare module '$env/dynamic/public' {
	export const env: {
		[key: `PUBLIC_${string}`]: string | undefined;
	}
}
