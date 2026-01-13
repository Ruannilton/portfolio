/**
 * app.js - Funções globais da aplicação
 */

// Navegação
function redirectToSearch() {
    window.location.href = '/app/search';
}

function login() {
    window.location.href = '/app/login';
}

function logout() {
    fetch('/auth/logout', { method: 'GET' })
        .then(response => {
            if (response.redirected) {
                window.location.href = response.url;
            }
        })
        .catch(error => console.error('Logout failed:', error));
}

// Controle de seções
function showSection(sectionId) {
    const sections = ['portfolio', 'conversations', 'stats'];
    sections.forEach(id => {
        const section = document.getElementById(id);
        if (section) {
            if (id === sectionId) {
                section.classList.remove('hidden');
            } else {
                section.classList.add('hidden');
            }
        }
    });
}
