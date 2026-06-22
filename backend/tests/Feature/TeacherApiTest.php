<?php

namespace Tests\Feature\Api;

use PHPUnit\Framework\TestCase;
use Illuminate\Support\Facades\Artisan;

class TeacherApiTest extends TestCase
{
    protected function setUp(): void
    {
        parent::setUp();
        Artisan::call('migrate:fresh');
    }

    /** @test */
    public function teacher_can_login_with_valid_credentials()
    {
        // Create test teacher
        $response = $this->postJson('/api/auth/teacher/login', [
            'username' => 'SCH001_T001',
            'password' => 'password123',
        ]);

        $response->assertStatus(200)
                 ->assertJsonStructure(['token', 'user']);
    }

    /** @test */
    public function teacher_cannot_login_with_invalid_credentials()
    {
        $response = $this->postJson('/api/auth/teacher/login', [
            'username' => 'invalid',
            'password' => 'wrong',
        ]);

        $response->assertStatus(401);
    }

    /** @test */
    public function teacher_can_view_student_list()
    {
        $token = $this->loginAsTeacher();

        $response = $this->withHeader('Authorization', 'Bearer ' . $token)
                         ->getJson('/api/teacher/students');

        $response->assertStatus(200)
                 ->assertJsonStructure(['data', 'total']);
    }

    /** @test */
    public function teacher_can_add_score_for_student()
    {
        $token = $this->loginAsTeacher();
        $studentId = 1;

        $response = $this->withHeader('Authorization', 'Bearer ' . $token)
                         ->postJson("/api/teacher/students/{$studentId}/add-score", [
                             'points' => 10,
                             'reason' => 'Test reason',
                         ]);

        $response->assertStatus(200);
        $this->assertDatabaseHas('scores', [
            'student_id' => $studentId,
            'change_value' => 10,
        ]);
    }

    protected function loginAsTeacher()
    {
        $response = $this->postJson('/api/auth/teacher/login', [
            'username' => 'SCH001_T001',
            'password' => 'password123',
        ]);

        return $response->json('token');
    }
}
