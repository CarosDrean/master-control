import { w as slot } from "../../chunks/index2.js";
import "@sveltejs/kit/internal";
import "../../chunks/exports.js";
import "../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../chunks/state.svelte.js";
function _layout($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    $$renderer2.push(`<header class="svelte-12qhfyh"><div class="left svelte-12qhfyh"><button class="title-btn svelte-12qhfyh" type="button">Master Control</button></div> <div class="right svelte-12qhfyh"><button class="icon-btn-btn svelte-12qhfyh" type="button"><img src="https://api.iconify.design/tabler/settings.svg" alt="settings" class="icon-btn svelte-12qhfyh"/></button> <button class="icon-btn-btn svelte-12qhfyh" type="button"><img src="https://api.iconify.design/tabler/layout-grid.svg" alt="pages" class="icon-btn svelte-12qhfyh"/></button></div></header> <main class="svelte-12qhfyh"><!--[-->`);
    slot($$renderer2, $$props, "default", {});
    $$renderer2.push(`<!--]--></main>`);
  });
}
export {
  _layout as default
};
