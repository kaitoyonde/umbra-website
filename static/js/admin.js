/* ==========================================================================
   Umbra CMS — Admin JavaScript
   Requires: admin.css
   ========================================================================== */

(function () {
  'use strict';

  // ========================================================================
  // UTILITY
  // ========================================================================

  const $ = (sel, ctx) => (ctx || document).querySelector(sel);
  const $$ = (sel, ctx) => Array.from((ctx || document).querySelectorAll(sel));

  function debounce(fn, ms) {
    let timer;
    return function (...args) {
      clearTimeout(timer);
      timer = setTimeout(() => fn.apply(this, args), ms);
    };
  }

  // ========================================================================
  // SIDEBAR
  // ========================================================================

  const Sidebar = {
    KEY: 'umbra_admin_sidebar_collapsed',

    init() {
      this.el = $('.admin-sidebar');
      this.content = $('.admin-content');
      this.topbar = $('.admin-topbar');
      this.toggleBtn = $('.admin-sidebar-toggle');

      if (!this.el) return;

      // Restore state
      const saved = localStorage.getItem(this.KEY);
      if (saved === 'true') {
        this.el.classList.add('collapsed');
        if (this.content) this.content.classList.add('sidebar-collapsed');
        if (this.topbar) this.topbar.classList.add('sidebar-collapsed');
      }

      this.toggleBtn?.addEventListener('click', () => this.toggle());

      // Nav item active state
      $$('.admin-nav-item').forEach(item => {
        if (item.getAttribute('href') === window.location.pathname) {
          item.classList.add('active');
        }
        // Expand parent subnav if child is active
        const sub = item.nextElementSibling;
        if (sub?.classList.contains('admin-nav-sub')) {
          const hasActive = sub.querySelector('.admin-nav-sub-item.active');
          if (hasActive) sub.classList.remove('collapsed');
        }
      });

      // Subnav toggle
      $$('.admin-nav-item.has-sub').forEach(item => {
        item.addEventListener('click', e => {
          e.preventDefault();
          const sub = item.nextElementSibling;
          if (sub?.classList.contains('admin-nav-sub')) {
            sub.classList.toggle('collapsed');
          }
        });
      });

      // Mobile: expand overlay on tablet+
      if (window.innerWidth <= 1024) {
        this.el.classList.add('collapsed');
        if (this.content) this.content.classList.add('sidebar-collapsed');
        if (this.topbar) this.topbar.classList.add('sidebar-collapsed');
      }

      // Hover to expand when collapsed on desktop
      if (window.innerWidth > 1024) {
        this.el.addEventListener('mouseenter', () => {
          if (this.el.classList.contains('collapsed')) {
            this.el.classList.add('expanded');
          }
        });
        this.el.addEventListener('mouseleave', () => {
          this.el.classList.remove('expanded');
        });
      }
    },

    toggle() {
      const collapsed = this.el.classList.toggle('collapsed');
      this.content?.classList.toggle('sidebar-collapsed', collapsed);
      this.topbar?.classList.toggle('sidebar-collapsed', collapsed);
      localStorage.setItem(this.KEY, collapsed);
    }
  };

  // ========================================================================
  // TOAST
  // ========================================================================

  const Toast = {
    container: null,

    init() {
      this.container = $('.toast-container');
      if (!this.container) {
        this.container = document.createElement('div');
        this.container.className = 'toast-container';
        document.body.appendChild(this.container);
      }
    },

    show(opts) {
      const { type = 'info', title, message, duration = 4000 } = opts;
      const icons = { success: '✓', error: '✕', warning: '!', info: 'i' };
      const el = document.createElement('div');
      el.className = `toast ${type}`;
      el.innerHTML = `
        <div class="toast-icon">${icons[type] || 'i'}</div>
        <div class="toast-body">
          <div class="toast-title">${this.esc(title)}</div>
          ${message ? `<div class="toast-message">${this.esc(message)}</div>` : ''}
        </div>
        <button class="toast-close" type="button">&times;</button>
      `;
      el.querySelector('.toast-close').addEventListener('click', () => this.dismiss(el));
      this.container.appendChild(el);

      // Trigger animation
      requestAnimationFrame(() => el.classList.add('show'));

      if (duration > 0) {
        this._timer = setTimeout(() => this.dismiss(el), duration);
        el.addEventListener('mouseenter', () => clearTimeout(this._timer));
        el.addEventListener('mouseleave', () => {
          this._timer = setTimeout(() => this.dismiss(el), duration);
        });
      }
    },

    dismiss(el) {
      el.classList.remove('show');
      el.classList.add('exit');
      setTimeout(() => el.remove(), 300);
    },

    esc(s) {
      const d = document.createElement('div');
      d.textContent = s;
      return d.innerHTML;
    }
  };

  // ========================================================================
  // MODAL
  // ========================================================================

  const Modal = {
    active: null,

    open(opts) {
      const { title, body, size = 'md', footer = '', onClose } = opts;
      this.close();

      const overlay = document.createElement('div');
      overlay.className = 'modal-overlay';
      overlay.innerHTML = `
        <div class="modal modal-${size}">
          <div class="modal-header">
            <h3 class="modal-title">${title}</h3>
            <button class="modal-close" type="button">&times;</button>
          </div>
          <div class="modal-body">${body}</div>
          ${footer ? `<div class="modal-footer">${footer}</div>` : ''}
        </div>
      `;

      overlay.querySelector('.modal-close').addEventListener('click', () => this.close());
      overlay.addEventListener('click', e => { if (e.target === overlay) this.close(); });
      document.addEventListener('keydown', this._keyHandler = e => { if (e.key === 'Escape') this.close(); });

      document.body.appendChild(overlay);
      requestAnimationFrame(() => overlay.classList.add('open'));
      this.active = overlay;
      if (onClose) this._onClose = onClose;
    },

    close() {
      if (!this.active) return;
      this.active.classList.remove('open');
      setTimeout(() => { this.active?.remove(); this.active = null; }, 200);
      if (this._keyHandler) document.removeEventListener('keydown', this._keyHandler);
      if (this._onClose) { this._onClose(); this._onClose = null; }
    }
  };

  // ========================================================================
  // CONFIRM DIALOG
  // ========================================================================

  const Confirm = {
    async show(opts) {
      const { title = 'Konfirmasi', message = 'Yakin ingin melanjutkan?', confirmText = 'Ya, hapus', cancelText = 'Batal', danger = true } = opts;

      return new Promise(resolve => {
        const footer = `
          <button class="btn btn-ghost" data-action="cancel">${cancelText}</button>
          <button class="btn ${danger ? 'btn-danger' : 'btn-primary'}" data-action="confirm">${confirmText}</button>
        `;
        Modal.open({
          title,
          size: 'sm',
          body: `<div class="modal-confirm">
            <div class="modal-icon">${danger ? '⚠' : 'ℹ'}</div>
            <div class="modal-message">${message}</div>
          </div>`,
          footer,
          onClose: () => resolve(false)
        });

        Modal.active.querySelector('[data-action="confirm"]').addEventListener('click', () => { Modal.close(); resolve(true); });
        Modal.active.querySelector('[data-action="cancel"]').addEventListener('click', () => { Modal.close(); resolve(false); });
      });
    }
  };

  // ========================================================================
  // IMAGE PICKER
  // ========================================================================

  const ImagePicker = {
    async open(opts) {
      const { images, selectedId, onSelect } = opts;
      let currentPage = 0;
      const perPage = 20;
      let filtered = images;

      const renderGrid = (filter) => {
        const q = (filter || '').toLowerCase();
        filtered = q ? images.filter(i => i.name.toLowerCase().includes(q) || i.alt.toLowerCase().includes(q)) : images;
        currentPage = 0;
        return renderPage();
      };

      const renderPage = () => {
        const start = currentPage * perPage;
        const page = filtered.slice(start, start + perPage);
        const totalPages = Math.ceil(filtered.length / perPage) || 1;

        if (!page.length) {
          return `<div class="image-picker-empty"><div class="image-picker-empty-icon">📷</div><div>Tidak ada gambar ditemukan</div></div>`;
        }

        let html = '<div class="image-picker-grid">';
        page.forEach(img => {
          const sel = img.id === selectedId ? ' selected' : '';
          html += `<div class="image-picker-card${sel}" data-id="${img.id}">
            <img src="${img.url}" alt="${img.alt}" loading="lazy" />
            <div class="image-picker-card-name">${img.name}</div>
          </div>`;
        });
        html += '</div>';

        if (totalPages > 1) {
          html += `<div class="image-picker-pagination">
            <button ${currentPage === 0 ? 'disabled' : ''} data-page="prev">← Sebelumnya</button>
            <span class="page-info">${currentPage + 1} / ${totalPages}</span>
            <button ${currentPage >= totalPages - 1 ? 'disabled' : ''} data-page="next">Selanjutnya →</button>
          </div>`;
        }
        return html;
      };

      const body = `
        <div class="image-picker-tabs">
          <button class="image-picker-tab active" data-tab="library">Pustaka Gambar</button>
          <button class="image-picker-tab" data-tab="upload">Upload Baru</button>
        </div>
        <div class="image-picker-tab-content active" data-tab-content="library">
          <div class="image-picker-search">
            <span class="search-icon">🔍</span>
            <input type="text" placeholder="Cari gambar..." id="ip-search" />
          </div>
          <div id="ip-grid">${renderGrid()}</div>
        </div>
        <div class="image-picker-tab-content" data-tab-content="upload">
          <div class="upload-dropzone" id="ip-dropzone">
            <div class="upload-dropzone-icon">📁</div>
            <div class="upload-dropzone-text">Seret gambar ke sini atau klik untuk memilih</div>
            <div class="upload-dropzone-hint">PNG, JPG, WebP — maks 10 MB</div>
            <div class="upload-dropzone-error" id="ip-dropzone-error"></div>
          </div>
          <div class="upload-preview" id="ip-preview">
            <img class="upload-preview-img" id="ip-preview-img" />
            <div class="upload-preview-fields">
              <input type="text" placeholder="Nama gambar" id="ip-upload-name" />
              <input type="text" placeholder="Alt text" id="ip-upload-alt" />
              <button class="btn btn-primary" id="ip-upload-btn">Upload</button>
            </div>
          </div>
          <div class="upload-progress" id="ip-upload-progress">
            <div class="upload-progress-bar"><div class="upload-progress-fill" id="ip-upload-fill"></div></div>
            <span class="upload-progress-text" id="ip-upload-text">Mengupload…</span>
          </div>
        </div>
      `;

      const footer = `<button class="btn btn-ghost" id="ip-cancel">Batal</button>
        <button class="btn btn-primary" id="ip-select" disabled>Pilih</button>`;

      Modal.open({
        title: 'Pilih Gambar',
        size: 'lg',
        body,
        footer
      });

      // Init tab switching
      $$('.image-picker-tab', Modal.active).forEach(tab => {
        tab.addEventListener('click', () => {
          $$('.image-picker-tab', Modal.active).forEach(t => t.classList.remove('active'));
          $$('[data-tab-content]', Modal.active).forEach(c => c.classList.remove('active'));
          tab.classList.add('active');
          const content = $(`[data-tab-content="${tab.dataset.tab}"]`, Modal.active);
          if (content) content.classList.add('active');
        });
      });

      // Init search
      const searchInput = $('#ip-search', Modal.active);
      if (searchInput) {
        searchInput.addEventListener('input', debounce(function () {
          $('#ip-grid', Modal.active).innerHTML = renderGrid(this.value);
          bindGridClicks();
        }, 300));
      }

      const bindGridClicks = () => {
        $$('.image-picker-card', Modal.active).forEach(card => {
          card.addEventListener('click', () => {
            $$('.image-picker-card', Modal.active).forEach(c => c.classList.remove('selected'));
            card.classList.add('selected');
            $('#ip-select', Modal.active).disabled = false;
            window._ipSelectedId = card.dataset.id;
          });
        });

        // Pagination
        $$('[data-page]', Modal.active).forEach(btn => {
          btn.addEventListener('click', () => {
            if (btn.dataset.page === 'prev' && currentPage > 0) currentPage--;
            if (btn.dataset.page === 'next' && (currentPage + 1) * perPage < filtered.length) currentPage++;
            $('#ip-grid', Modal.active).innerHTML = renderPage();
            bindGridClicks();
          });
        });
      };
      bindGridClicks();

      // Init upload
      const dropzone = $('#ip-dropzone', Modal.active);
      const fileInput = document.createElement('input');
      fileInput.type = 'file';
      fileInput.accept = 'image/*';
      fileInput.style.display = 'none';
      Modal.active.appendChild(fileInput);

      dropzone.addEventListener('click', () => fileInput.click());

      dropzone.addEventListener('dragover', e => { e.preventDefault(); dropzone.classList.add('dragover'); });
      dropzone.addEventListener('dragleave', () => dropzone.classList.remove('dragover'));
      dropzone.addEventListener('drop', e => {
        e.preventDefault();
        dropzone.classList.remove('dragover');
        if (e.dataTransfer.files.length) handleFile(e.dataTransfer.files[0]);
      });
      fileInput.addEventListener('change', () => { if (fileInput.files.length) handleFile(fileInput.files[0]); });

      const handleFile = (file) => {
        if (file.size > 10 * 1024 * 1024) {
          $('#ip-dropzone-error', Modal.active).textContent = 'File terlalu besar. Maksimal 10 MB.';
          $('#ip-dropzone-error', Modal.active).style.display = 'block';
          return;
        }
        $('#ip-dropzone-error', Modal.active).style.display = 'none';
        const reader = new FileReader();
        reader.onload = e => {
          $('#ip-preview-img', Modal.active).src = e.target.result;
          $('#ip-preview', Modal.active).classList.add('show');
          window._ipUploadFile = file;
        };
        reader.readAsDataURL(file);
      };

      $('#ip-upload-btn', Modal.active)?.addEventListener('click', async () => {
        const name = $('#ip-upload-name', Modal.active)?.value;
        const alt = $('#ip-upload-alt', Modal.active)?.value;
        if (!name || !alt || !window._ipUploadFile) {
          Toast.show({ type: 'warning', title: 'Lengkapi data', message: 'Isi nama, alt text, dan pilih file.' });
          return;
        }
        const formData = new FormData();
        formData.append('name', name);
        formData.append('alt', alt);
        formData.append('description', '');
        formData.append('file', window._ipUploadFile);

        $('#ip-upload-progress', Modal.active).classList.add('show');

        try {
          const res = await fetch('/admin/images', { method: 'POST', body: formData });
          if (!res.ok) throw new Error('Upload gagal');
          Toast.show({ type: 'success', title: 'Gambar terupload', message: `${name} berhasil diupload.` });
          Modal.close();
          onSelect(null, true); // signal refresh
        } catch (err) {
          Toast.show({ type: 'error', title: 'Upload gagal', message: err.message });
          $('#ip-upload-progress', Modal.active).classList.remove('show');
        }
      });

      // Select button
      $('#ip-select', Modal.active).addEventListener('click', () => {
        const id = window._ipSelectedId;
        if (id) {
          const img = images.find(i => i.id === id);
          Modal.close();
          onSelect(img, false);
        }
      });

      $('#ip-cancel', Modal.active).addEventListener('click', () => Modal.close());
    }
  };

  // ========================================================================
  // LIGHTBOX
  // ========================================================================

  const Lightbox = {
    init() {
      document.addEventListener('click', e => {
        const img = e.target.closest('[data-lightbox]');
        if (img) {
          e.preventDefault();
          this.open(img.getAttribute('src') || img.href, img.getAttribute('alt') || '');
        }
      });
    },

    open(src, alt) {
      const overlay = document.createElement('div');
      overlay.className = 'lightbox open';
      overlay.innerHTML = `<button class="lightbox-close" type="button">&times;</button><img src="${src}" alt="${alt}" />`;
      overlay.querySelector('.lightbox-close').addEventListener('click', () => this.close(overlay));
      overlay.addEventListener('click', e => { if (e.target === overlay) this.close(overlay); });
      document.addEventListener('keydown', this._lbKeyHandler = e => { if (e.key === 'Escape') this.close(overlay); });
      document.body.appendChild(overlay);
    },

    close(overlay) {
      overlay.classList.remove('open');
      setTimeout(() => overlay.remove(), 200);
      if (this._lbKeyHandler) document.removeEventListener('keydown', this._lbKeyHandler);
    }
  };

  // ========================================================================
  // TABS
  // ========================================================================

  const Tabs = {
    init(container) {
      container = container || document;
      $$('.admin-tabs', container).forEach(tabs => {
        const btns = $$('.admin-tab', tabs);
        const contents = $$('.admin-tab-content', container);
        btns.forEach(btn => {
          btn.addEventListener('click', () => {
            btns.forEach(b => b.classList.remove('active'));
            contents.forEach(c => c.classList.remove('active'));
            btn.classList.add('active');
            const target = $(`#${btn.dataset.tab}`, container) || $(`[data-tab-content="${btn.dataset.tab}"]`, container);
            if (target) target.classList.add('active');
          });
        });
      });
    }
  };

  // ========================================================================
  // CONFIRM DELETE (replaces native confirm())
  // ========================================================================

  function initConfirmDeletes() {
    document.addEventListener('click', async e => {
      const link = e.target.closest('[data-confirm]');
      if (!link) return;
      e.preventDefault();
      const msg = link.dataset.confirm || 'Yakin ingin menghapus?';
      const confirmed = await Confirm.show({ message: msg });
      if (confirmed) {
        window.location.href = link.getAttribute('href');
      }
    });
  }

  // ========================================================================
  // TOAST FLASH MESSAGES
  // ========================================================================

  function initFlashToasts() {
    // Look for a meta tag with flash data
    const flash = document.querySelector('meta[name="flash"]');
    if (flash) {
      try {
        const data = JSON.parse(flash.content);
        if (data.message) {
          Toast.show({ type: data.type || 'info', title: data.title || '', message: data.message });
        }
      } catch (e) { /* ignore */ }
    }
    // Also check URL params
    const params = new URLSearchParams(window.location.search);
    const flashMsg = params.get('flash');
    if (flashMsg) {
      Toast.show({ type: params.get('flash_type') || 'success', title: params.get('flash_title') || '', message: decodeURIComponent(flashMsg) });
    }
  }

  // ========================================================================
  // SIDEBAR DIVISION LINKS — dynamic from DOM
  // ========================================================================

  function initDivisionNav() {
    // Division nav items are rendered server-side in admin-base.html
    // Just activate the current one
    $$('.admin-nav-sub-item').forEach(item => {
      if (item.getAttribute('href') === window.location.pathname) {
        item.classList.add('active');
      }
    });
  }

  // ========================================================================
  // INIT
  // ========================================================================

  document.addEventListener('DOMContentLoaded', () => {
    Sidebar.init();
    Toast.init();
    Lightbox.init();
    Tabs.init();
    initConfirmDeletes();
    initFlashToasts();
    initDivisionNav();
  });

  // ========================================================================
  // EXPOSE GLOBALS (for inline use in templates)
  // ========================================================================

  window.Admin = {
    Toast,
    Modal,
    Confirm,
    ImagePicker,
    Lightbox,
    Tabs,
    Sidebar
  };

})();
