htmx.config.globalViewTransitions = true
document.body.addEventListener('htmx:beforeTransition', (e) => {
  console.log(e);
  if (e.target !== document.body) {
    e.preventDefault();
  }
});
