import { y as attr, x as ensure_array_like } from "../../../chunks/index2.js";
import "@sveltejs/kit/internal";
import "../../../chunks/exports.js";
import "../../../chunks/utils.js";
import "@sveltejs/kit/internal/server";
import "../../../chunks/state.svelte.js";
import { k as escape_html } from "../../../chunks/context.js";
function _page($$renderer, $$props) {
  $$renderer.component(($$renderer2) => {
    let pages = [];
    let id = "";
    let name = "";
    let icon = "";
    let editing = null;
    let editName = "";
    let editIcon = "";
    $$renderer2.push(`<div class="container svelte-1bueben"><button class="btn-back svelte-1bueben">← Volver</button> <h1 class="svelte-1bueben">Páginas</h1> <div class="panel svelte-1bueben"><div class="field svelte-1bueben"><label for="id" class="svelte-1bueben">ID</label> <input id="id"${attr("value", id)} placeholder="general" class="svelte-1bueben"/></div> <div class="field svelte-1bueben"><label for="name" class="svelte-1bueben">Nombre</label> <input id="name"${attr("value", name)} placeholder="General" class="svelte-1bueben"/></div> <div class="field svelte-1bueben"><label for="icon" class="svelte-1bueben">Icono (emoji)</label> <input id="icon"${attr("value", icon)} placeholder="✨" class="svelte-1bueben"/></div> <button class="btn svelte-1bueben">Añadir página</button></div> <div class="panel svelte-1bueben">`);
    if (pages.length === 0) {
      $$renderer2.push("<!--[-->");
      $$renderer2.push(`<div>No hay páginas</div>`);
    } else {
      $$renderer2.push("<!--[!-->");
      $$renderer2.push(`<!--[-->`);
      const each_array = ensure_array_like(pages);
      for (let $$index = 0, $$length = each_array.length; $$index < $$length; $$index++) {
        let p = each_array[$$index];
        $$renderer2.push(`<div class="row svelte-1bueben"><span class="page-icon svelte-1bueben">${escape_html(p.icon)}</span> `);
        if (editing === p.id) {
          $$renderer2.push("<!--[-->");
          $$renderer2.push(`<input${attr("value", editName)} placeholder="Nombre" class="svelte-1bueben"/> <input${attr("value", editIcon)} placeholder="Icono" style="width:80px" class="svelte-1bueben"/> <button class="btn svelte-1bueben">Guardar</button>`);
        } else {
          $$renderer2.push("<!--[!-->");
          $$renderer2.push(`<span class="page-name svelte-1bueben">${escape_html(p.name)}</span> <button class="btn svelte-1bueben">Editar</button> <button class="btn svelte-1bueben">Eliminar</button>`);
        }
        $$renderer2.push(`<!--]--></div>`);
      }
      $$renderer2.push(`<!--]-->`);
    }
    $$renderer2.push(`<!--]--></div></div>`);
  });
}
export {
  _page as default
};
