import { y as attr, x as ensure_array_like } from "../../../../chunks/index2.js";
import "@sveltejs/kit/internal";
import "../../../../chunks/exports.js";
import "../../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../../chunks/state.svelte.js";
import { a as appIconUrl } from "../../../../chunks/HotkeyRecorder.svelte_svelte_type_style_lang.js";
import { k as escape_html } from "../../../../chunks/context.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let selectedApp, isImage, isAppType;
    let actionTypes = [];
    let pages = [];
    let apps = [];
    let type = "";
    let page = "general";
    let id = "";
    let label = "";
    let command = "";
    let icon = "";
    selectedApp = apps.find((a) => a.name === command);
    isImage = icon?.startsWith("http") || icon?.startsWith("/");
    isAppType = type === "app";
    let $$settled = true;
    let $$inner_renderer;
    function $$render_inner($$renderer3) {
      $$renderer3.push(`<div class="container svelte-1mq7x49"><button class="btn-back svelte-1mq7x49">← Regresar</button> <h1 class="svelte-1mq7x49">Nueva acción</h1> <div class="preview svelte-1mq7x49">`);
      if (isImage) {
        $$renderer3.push("<!--[-->");
        $$renderer3.push(`<img${attr("src", icon)} alt="" class="svelte-1mq7x49"/>`);
      } else {
        $$renderer3.push("<!--[!-->");
        $$renderer3.push(`<div class="emoji svelte-1mq7x49">${escape_html(icon)}</div>`);
      }
      $$renderer3.push(`<!--]--> <div class="preview-label svelte-1mq7x49">${escape_html(label)}</div></div> `);
      {
        $$renderer3.push("<!--[!-->");
      }
      $$renderer3.push(`<!--]--> <div class="panel svelte-1mq7x49"><div class="field svelte-1mq7x49"><label for="id" class="svelte-1mq7x49">ID</label> <input id="id"${attr("value", id)} placeholder="se genera solo desde el label" class="svelte-1mq7x49"/></div> <div class="field svelte-1mq7x49"><label for="label" class="svelte-1mq7x49">Label</label> <input id="label"${attr("value", label)} placeholder="Nombre visible" class="svelte-1mq7x49"/></div> <div class="field svelte-1mq7x49"><label for="icon" class="svelte-1mq7x49">Icono (emoji o URL)</label> <input id="icon"${attr("value", icon)} placeholder="😀 o https://…" class="svelte-1mq7x49"/></div> <div class="field svelte-1mq7x49"><label for="type" class="svelte-1mq7x49">Tipo</label> `);
      $$renderer3.select(
        { id: "type", value: type, class: "" },
        ($$renderer4) => {
          $$renderer4.push(`<!--[-->`);
          const each_array = ensure_array_like(actionTypes);
          for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
            let t = each_array[$$index];
            $$renderer4.option({ value: t.id }, ($$renderer5) => {
              $$renderer5.push(`${escape_html(t.name)} — ${escape_html(t.description)}`);
            });
          }
          $$renderer4.push(`<!--]-->`);
        },
        "svelte-1mq7x49"
      );
      $$renderer3.push(`</div> <div class="field svelte-1mq7x49"><label for="page" class="svelte-1mq7x49">Página</label> `);
      $$renderer3.select(
        { id: "page", value: page, class: "" },
        ($$renderer4) => {
          $$renderer4.push(`<!--[-->`);
          const each_array_1 = ensure_array_like(pages);
          for (let $$index_1 = 0, $$length = each_array_1.length; $$index_1 < $$length; $$index_1++) {
            let p = each_array_1[$$index_1];
            $$renderer4.option({ value: p.id }, ($$renderer5) => {
              $$renderer5.push(`${escape_html(p.icon)} ${escape_html(p.name)}`);
            });
          }
          $$renderer4.push(`<!--]-->`);
        },
        "svelte-1mq7x49"
      );
      $$renderer3.push(`</div> `);
      if (isAppType) {
        $$renderer3.push("<!--[-->");
        $$renderer3.push(`<div class="field svelte-1mq7x49"><label for="app-search" class="svelte-1mq7x49">App</label> `);
        {
          $$renderer3.push("<!--[!-->");
          $$renderer3.push(`<button type="button" class="app-select svelte-1mq7x49">`);
          if (selectedApp) {
            $$renderer3.push("<!--[-->");
            $$renderer3.push(`<img${attr("src", appIconUrl(selectedApp.path))}${attr("alt", selectedApp.name)} class="svelte-1mq7x49"/> <span>${escape_html(selectedApp.name)}</span>`);
          } else {
            $$renderer3.push("<!--[!-->");
            $$renderer3.push(`<span class="placeholder svelte-1mq7x49">Selecciona una app…</span>`);
          }
          $$renderer3.push(`<!--]--> <span class="chevron svelte-1mq7x49">▾</span></button>`);
        }
        $$renderer3.push(`<!--]--></div>`);
      } else {
        $$renderer3.push("<!--[!-->");
        {
          $$renderer3.push("<!--[!-->");
          $$renderer3.push(`<div class="field svelte-1mq7x49"><label for="command" class="svelte-1mq7x49">Command / Script</label> <textarea id="command" placeholder="osascript -e 'tell application …'" class="svelte-1mq7x49">`);
          const $$body = escape_html(command);
          if ($$body) {
            $$renderer3.push(`${$$body}`);
          }
          $$renderer3.push(`</textarea></div>`);
        }
        $$renderer3.push(`<!--]-->`);
      }
      $$renderer3.push(`<!--]--> <button class="btn-save svelte-1mq7x49">Guardar</button></div></div>`);
    }
    do {
      $$settled = true;
      $$inner_renderer = $$renderer2.copy();
      $$render_inner($$inner_renderer);
    } while (!$$settled);
    $$renderer2.subsume($$inner_renderer);
  });
}
export {
  _page as default
};
