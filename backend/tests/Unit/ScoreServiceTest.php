<?php

namespace Tests\Unit;

use PHPUnit\Framework\TestCase;
use App\Services\ScoreService;
use App\Models\Student;
use App\Models\Score;

class ScoreServiceTest extends TestCase
{
    protected function setUp(): void
    {
        parent::setUp();
        // Mock database or use in-memory SQLite
    }

    /** @test */
    public function it_can_add_score_to_student()
    {
        // Arrange
        $student = Student::create([
            'name' => 'Test Student',
            'class_id' => 1,
            'total_score' => 0,
        ]);

        // Act
        $result = ScoreService::addScore($student->id, 10, 'Test reason');

        // Assert
        $this->assertTrue($result);
        $this->assertEquals(10, $student->fresh()->total_score);
    }

    /** @test */
    public function it_can_calculate_level_from_experience()
    {
        // Test level calculation logic
        $this->assertEquals(1, ScoreService::calculateLevel(0));
        $this->assertEquals(1, ScoreService::calculateLevel(99));
        $this->assertEquals(2, ScoreService::calculateLevel(100));
        $this->assertEquals(3, ScoreService::calculateLevel(300));
    }

    /** @test */
    public function it_records_score_history()
    {
        $student = Student::create([
            'name' => 'Test Student',
            'class_id' => 1,
            'total_score' => 0,
        ]);

        ScoreService::addScore($student->id, 10, 'Test reason');

        $this->assertDatabaseHas('scores', [
            'student_id' => $student->id,
            'change_value' => 10,
            'reason' => 'Test reason',
        ]);
    }
}
