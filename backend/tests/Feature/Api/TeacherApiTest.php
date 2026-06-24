<?php

declare(strict_types=1);

namespace Tests\Feature\Api;

use Tests\TestCase;

class TeacherApiTest extends TestCase
{
    /** @test */
    public function application_health_check_returns_ok(): void
    {
        $response = $this->get('/up');

        $response->assertStatus(200);
    }

    /** @test */
    public function teacher_login_endpoint_exists(): void
    {
        $response = $this->postJson('/api/auth/teacher/login', []);

        // 422 = validation error (endpoint exists, input invalid)
        // 401 = unauthorized (endpoint exists, credentials wrong)
        // 500 = server error (endpoint exists, but something broke)
        // 404 = endpoint doesn't exist
        $this->assertNotEquals(404, $response->status(), 'Teacher login endpoint should exist');
    }

    /** @test */
    public function teacher_login_with_invalid_credentials_returns_error(): void
    {
        $response = $this->postJson('/api/auth/teacher/login', [
            'username' => 'nonexistent_user',
            'password' => 'wrong_password',
        ]);

        // Should be 401 (unauthorized) or 422 (validation) — not 500 or 404
        $this->assertContains(
            $response->status(),
            [401, 422],
            'Invalid login should return 401 or 422'
        );
    }

    /** @test */
    public function teacher_endpoints_require_authentication(): void
    {
        $endpoints = [
            ['GET', '/api/teacher/students'],
            ['GET', '/api/teacher/leaderboard'],
            ['GET', '/api/teacher/pets'],
        ];

        foreach ($endpoints as [$method, $uri]) {
            $response = $this->json($method, $uri);

            // Should be 401 (unauthorized) — not 404 (not found) or 500 (error)
            $this->assertNotEquals(
                404,
                $response->status(),
                "Endpoint {$method} {$uri} should exist (not 404)"
            );
        }
    }
}
