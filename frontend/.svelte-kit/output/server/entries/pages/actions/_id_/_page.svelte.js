import { z as bind_props, y as attr, x as ensure_array_like } from "../../../../chunks/index2.js";
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
    let params = $$props["params"];
    let actionTypes = [];
    let pages = [];
    let apps = [];
    let label = "";
    let type = "";
    let command = "";
    let icon = "";
    let page = "general";
    selectedApp = apps.find((a) => a.name === command);
    isImage = icon?.startsWith("http") || icon?.startsWith("/");
    isAppType = type === "app";
    let $$settled = true;
    let $$inner_renderer;
    function $$render_inner($$renderer3) {
      $$renderer3.push(`<div class="container svelte-3yf9ho"><button class="btn-back svelte-3yf9ho">← Regresar</button> <h1 class="svelte-3yf9ho">Editar acción</h1> <div class="preview svelte-3yf9ho">`);
      if (isImage) {
        $$renderer3.push("<!--[-->");
        $$renderer3.push(`<img${attr("src", icon)} alt="" class="svelte-3yf9ho"/>`);
      } else {
        $$renderer3.push("<!--[!-->");
        $$renderer3.push(`<div class="emoji svelte-3yf9ho">${escape_html(icon)}</div>`);
      }
      $$renderer3.push(`<!--]--> <div class="preview-label svelte-3yf9ho">${escape_html(label)}</div></div> `);
      {
        $$renderer3.push("<!--[!-->");
      }
      $$renderer3.push(`<!--]--> <div class="panel svelte-3yf9ho"><div class="field svelte-3yf9ho"><label for="label-edit" class="svelte-3yf9ho">Label</label> <input id="label-edit"${attr("value", label)} placeholder="Nombre visible" class="svelte-3yf9ho"/></div> <div class="field svelte-3yf9ho"><label for="icon-edit" class="svelte-3yf9ho">Icono (emoji o URL)</label> <input id="icon-edit"${attr("value", icon)} placeholder="😀 o https://…" class="svelte-3yf9ho"/></div> <div class="field svelte-3yf9ho"><label for="type-edit" class="svelte-3yf9ho">Tipo</label> `);
      $$renderer3.select(
        { id: "type-edit", value: type, class: "" },
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
        "svelte-3yf9ho"
      );
      $$renderer3.push(`</div> <div class="field svelte-3yf9ho"><label for="page-edit" class="svelte-3yf9ho">Página</label> `);
      $$renderer3.select(
        { id: "page-edit", value: page, class: "" },
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
        "svelte-3yf9ho"
      );
      $$renderer3.push(`</div> `);
      if (isAppType) {
        $$renderer3.push("<!--[-->");
        $$renderer3.push(`<div class="field svelte-3yf9ho"><label for="app-search" class="svelte-3yf9ho">App</label> `);
        {
          $$renderer3.push("<!--[!-->");
          $$renderer3.push(`<button type="button" class="app-select svelte-3yf9ho">`);
          if (selectedApp) {
            $$renderer3.push("<!--[-->");
            $$renderer3.push(`<img${attr("src", appIconUrl(selectedApp.path))}${attr("alt", selectedApp.name)} class="svelte-3yf9ho"/> <span>${escape_html(selectedApp.name)}</span>`);
          } else {
            $$renderer3.push("<!--[!-->");
            $$renderer3.push(`<span class="placeholder svelte-3yf9ho">Selecciona una app…</span>`);
          }
          $$renderer3.push(`<!--]--> <span class="chevron svelte-3yf9ho">▾</span></button>`);
        }
        $$renderer3.push(`<!--]--></div>`);
      } else {
        $$renderer3.push("<!--[!-->");
        {
          $$renderer3.push("<!--[!-->");
          $$renderer3.push(`<div class="field svelte-3yf9ho"><label for="command-edit" class="svelte-3yf9ho">Command / Script</label> <textarea id="command-edit" class="svelte-3yf9ho">`);
          const $$body = escape_html(command);
          if ($$body) {
            $$renderer3.push(`${$$body}`);
          }
          $$renderer3.push(`</textarea></div>`);
        }
        $$renderer3.push(`<!--]-->`);
      }
      $$renderer3.push(`<!--]--> <button class="btn-save svelte-3yf9ho">▶ Probar</button> <button class="btn-save svelte-3yf9ho">Guardar cambios</button></div></div>`);
    }
    do {
      $$settled = true;
      $$inner_renderer = $$renderer2.copy();
      $$render_inner($$inner_renderer);
    } while (!$$settled);
    $$renderer2.subsume($$inner_renderer);
    bind_props($$props, { params });
  });
}
export {
  _page as default
};
