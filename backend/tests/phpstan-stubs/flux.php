<?php

declare(strict_types=1);

// Stub file for Livewire Flux package
// Tells PHPStan about the Flux class and the global flux() helper function

namespace Livewire\Flux;

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
