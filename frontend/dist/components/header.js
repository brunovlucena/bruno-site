// Header functionality
document.addEventListener('DOMContentLoaded', function() {
    const mobileMenuToggle = document.querySelector('.mobile-menu-toggle');
    const headerNav = document.querySelector('.header-nav');
    const searchToggle = document.querySelector('.search-toggle');
    const dropdownToggle = document.querySelector('.dropdown-toggle');
    const dropdownMenu = document.querySelector('.dropdown-menu');
    
    // Mobile menu toggle
    if (mobileMenuToggle && headerNav) {
        mobileMenuToggle.addEventListener('click', function() {
            const isOpen = headerNav.classList.contains('nav-open');
            
            if (isOpen) {
                headerNav.classList.remove('nav-open');
                mobileMenuToggle.classList.remove('active');
            } else {
                headerNav.classList.add('nav-open');
                mobileMenuToggle.classList.add('active');
            }
        });
    }
    
    // Dropdown toggle functionality
    if (dropdownToggle && dropdownMenu) {
        dropdownToggle.addEventListener('click', function(e) {
            e.stopPropagation();
            const isOpen = dropdownMenu.classList.contains('show');
            
            if (isOpen) {
                dropdownMenu.classList.remove('show');
                dropdownToggle.classList.remove('active');
                dropdownToggle.setAttribute('aria-expanded', 'false');
            } else {
                dropdownMenu.classList.add('show');
                dropdownToggle.classList.add('active');
                dropdownToggle.setAttribute('aria-expanded', 'true');
            }
        });
        
        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (!dropdownToggle.contains(e.target) && !dropdownMenu.contains(e.target)) {
                dropdownMenu.classList.remove('show');
                dropdownToggle.classList.remove('active');
                dropdownToggle.setAttribute('aria-expanded', 'false');
            }
        });
        
        // Close dropdown on escape key
        document.addEventListener('keydown', function(e) {
            if (e.key === 'Escape' && dropdownMenu.classList.contains('show')) {
                dropdownMenu.classList.remove('show');
                dropdownToggle.classList.remove('active');
                dropdownToggle.setAttribute('aria-expanded', 'false');
            }
        });
    }
    
    // Search toggle functionality
    if (searchToggle) {
        searchToggle.addEventListener('click', function() {
            // For now, just show an alert. This can be expanded later
            alert('Search functionality coming soon!');
        });
    }
    
    // Close mobile menu when clicking outside
    document.addEventListener('click', function(event) {
        if (headerNav && headerNav.classList.contains('nav-open')) {
            if (!headerNav.contains(event.target) && !mobileMenuToggle.contains(event.target)) {
                headerNav.classList.remove('nav-open');
                mobileMenuToggle.classList.remove('active');
            }
        }
    });
    
    // Keyboard navigation support
    document.addEventListener('keydown', function(event) {
        if (event.key === 'Escape' && headerNav && headerNav.classList.contains('nav-open')) {
            headerNav.classList.remove('nav-open');
            mobileMenuToggle.classList.remove('active');
        }
    });
    
    // Smooth scrolling for navigation links
    const navLinks = document.querySelectorAll('.nav-menu a[href^="#"]');
    navLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href');
            const targetElement = document.querySelector(targetId);
            
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
                
                // Close mobile menu after clicking a link
                if (headerNav && headerNav.classList.contains('nav-open')) {
                    headerNav.classList.remove('nav-open');
                    mobileMenuToggle.classList.remove('active');
                }
            }
        });
    });
}); 