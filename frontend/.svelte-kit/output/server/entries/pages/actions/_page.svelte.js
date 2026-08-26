import { x as ensure_array_like, y as attr } from "../../../chunks/index2.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/state.svelte.js";
import { k as escape_html } from "../../../chunks/context.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let actions = [];
    const isImage = (icon) => !!icon && (icon.startsWith("http") || icon.startsWith("/") || /\.(png|jpe?g|svg|webp|gif|ico)$/i.test(icon));
    $$renderer2.push(`<div class="container svelte-11zl00k"><div class="header-row svelte-11zl00k"><div class="title svelte-11zl00k">Administrar acciones</div> <button class="btn-create svelte-11zl00k">Nueva acción</button></div> <div class="panel svelte-11zl00k"><table class="svelte-11zl00k"><thead><tr><th class="svelte-11zl00k">ID</th><th class="svelte-11zl00k">Icono</th><th class="svelte-11zl00k">Label</th><th class="svelte-11zl00k">Tipo</th><th class="svelte-11zl00k">Página</th><th style="text-align:right" class="svelte-11zl00k">Acciones</th></tr></thead><tbody><!--[-->`);
    const each_array = ensure_array_like(actions);
    for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
      let a = each_array[$$index];
      $$renderer2.push(`<tr class="svelte-11zl00k"><td class="svelte-11zl00k">${escape_html(a.id)}</td><td class="svelte-11zl00k">`);
      if (isImage(a.icon)) {
        $$renderer2.push("<!--[-->");
        $$renderer2.push(`<img class="row-icon svelte-11zl00k"${attr("src", a.icon)}${attr("alt", a.label)}/>`);
      } else {
        $$renderer2.push("<!--[!-->");
        $$renderer2.push(`<span class="row-emoji svelte-11zl00k">${escape_html(a.icon)}</span>`);
      }
      $$renderer2.push(`<!--]--></td><td class="svelte-11zl00k">${escape_html(a.label)}</td><td class="svelte-11zl00k">${escape_html(a.type)}</td><td class="svelte-11zl00k">${escape_html(a.page ?? "general")}</td><td style="text-align:right" class="svelte-11zl00k"><button class="btn-small svelte-11zl00k">Editar</button> <button class="btn-small svelte-11zl00k">Eliminar</button></td></tr>`);
    }
    $$renderer2.push(`<!--]--></tbody></table></div></div>`);
  });
}
export {
  _page as default
};
