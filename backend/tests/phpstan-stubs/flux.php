<?php

/**
 * PHPStan stub for Livewire Flux helper function
 *
 * This file tells PHPStan about the flux() helper function
 * provided by the Livewire Flux package.
 */

declare(strict_types=1);

// Define the flux() function in the global namespace
if (!function_exists('flux')) {
    /**
     * Get the Flux instance for toast notifications, modals, etc.
     *
     * @return \Livewire\Flux\Flux
     */
    function flux()
    {
        // Stub implementation - returns a Flux instance
        return new \Livewire\Flux\Flux();
    }
}

namespace Livewire\Flux {
    /**
     * Flux class for type hints
     */
    class Flux
    {
        /**
         * Show a toast notification
         *
         * @param string $message
         * @param string $type
         * @return void
         */
        public function toast(string $message, string $type = 'info'): void
        {
        }

        /**
         * Show a modal
         *
         * @param string $name
         * @return void
         */
        public function modal(string $name): void
        {
        }
    }
}
