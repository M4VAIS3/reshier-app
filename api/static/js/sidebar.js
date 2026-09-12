/**
 * sidebar.js — Collapsible sidebar controller
 * Handles hamburger toggle, overlay, keyboard shortcut, and responsive behavior.
 */
(function () {
  'use strict';

  const OPEN_CLASS   = 'open';
  const PUSHED_CLASS = 'sidebar-pushed';
  const ACTIVE_CLASS = 'active';

  let sidebar     = null;
  let mainContent = null;
  let overlay     = null;
  let toggleBtn   = null;
  let isOpen      = false;

  /* ── helpers ── */
  function openSidebar() {
    isOpen = true;
    sidebar.classList.add(OPEN_CLASS);
    overlay.classList.add(ACTIVE_CLASS);
    toggleBtn.setAttribute('aria-expanded', 'true');
    // Only push content on larger screens
    if (window.innerWidth >= 768) {
      mainContent.classList.add(PUSHED_CLASS);
    }
    document.body.style.overflow = window.innerWidth < 768 ? 'hidden' : '';
  }

  function closeSidebar() {
    isOpen = false;
    sidebar.classList.remove(OPEN_CLASS);
    overlay.classList.remove(ACTIVE_CLASS);
    mainContent.classList.remove(PUSHED_CLASS);
    toggleBtn.setAttribute('aria-expanded', 'false');
    document.body.style.overflow = '';
  }

  function toggleSidebar() {
    if (isOpen) { closeSidebar(); } else { openSidebar(); }
  }

  /* ── init ── */
  function init() {
    sidebar     = document.getElementById('sidebar');
    mainContent = document.querySelector('.main-content');
    toggleBtn   = document.getElementById('sidebar-toggle');

    if (!sidebar || !mainContent || !toggleBtn) return;

    /* Create overlay if not already present */
    overlay = document.getElementById('sidebar-overlay');
    if (!overlay) {
      overlay = document.createElement('div');
      overlay.id = 'sidebar-overlay';
      overlay.className = 'sidebar-overlay';
      document.body.appendChild(overlay);
    }

    /* Button click */
    toggleBtn.addEventListener('click', toggleSidebar);

    /* Overlay click — close */
    overlay.addEventListener('click', closeSidebar);

    /* Keyboard: Escape to close */
    document.addEventListener('keydown', function (e) {
      if (e.key === 'Escape' && isOpen) closeSidebar();
    });

    /* On resize — if screen widens, sidebar stays closed by default */
    window.addEventListener('resize', function () {
      if (isOpen && window.innerWidth >= 768) {
        mainContent.classList.add(PUSHED_CLASS);
        document.body.style.overflow = '';
      } else if (isOpen && window.innerWidth < 768) {
        mainContent.classList.remove(PUSHED_CLASS);
        document.body.style.overflow = 'hidden';
      }
    });
  }

  /* ── run after DOM ready ── */
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
