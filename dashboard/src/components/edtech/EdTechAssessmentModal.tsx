import { useState, useEffect, useCallback } from 'react';
import { Sparkles, ChevronLeft, ChevronRight, Check, Loader2 } from 'lucide-react';
import { Modal, ModalHeader, ModalBody, ModalFooter } from '@/components/ui/modal';
import { Button } from '@/components/ui/button';
import { AssessmentStepper } from './AssessmentStepper';
import { InfrastructureStep } from './InfrastructureStep';
import { PainPointsStep } from './PainPointsStep';
import { GoalsStep } from './GoalsStep';
import { AIFollowUpStep } from './AIFollowUpStep';
import {
  useEdTechOptions,
  useEdTechProfile,
  useSaveEdTechProfile,
  useGenerateAI,
  useSubmitFollowUp,
  useCompleteProfile,
} from '@/hooks/useEdTechProfile';
import type {
  EdTechProfile,
  InfrastructureStepData,
  PainPointsStepData,
  GoalsStepData,
  DeviceTypes,
} from '@/types/edtech';

interface EdTechAssessmentModalProps {
  open: boolean;
  onClose: () => void;
  schoolId: string;
}

const defaultDeviceTypes: DeviceTypes = {
  laptops: 0,
  chromebooks: 0,
  tablets: 0,
  desktops: 0,
  other: 0,
};

const defaultInfrastructure: InfrastructureStepData = {
  totalDevices: 0,
  deviceTypes: defaultDeviceTypes,
  networkQuality: '',
  internetSpeed: '',
  lmsPlatform: '',
  existingSoftware: [],
  itStaffCount: 0,
  deviceAge: '',
};

const defaultPainPoints: PainPointsStepData = {
  painPoints: [],
  supportSatisfaction: 0,
  biggestChallenges: [],
  supportFrequency: '',
  avgResolutionTime: '',
  biggestFrustration: '',
  wishList: '',
};

const defaultGoals: GoalsStepData = {
  strategicGoals: [],
  budgetRange: '',
  timeline: '',
  expansionPlans: '',
  priorityRanking: [],
  decisionMakers: [],
};

// Question counts for each step
const STEP_QUESTION_COUNTS: Record<number, number> = {
  1: 7, // Infrastructure
  2: 7, // Pain Points
  3: 6, // Goals
  4: 1, // AI Follow-up (single view)
};

export function EdTechAssessmentModal({ open, onClose, schoolId }: EdTechAssessmentModalProps) {
  const [currentStep, setCurrentStep] = useState(1);
  const [currentQuestion, setCurrentQuestion] = useState(1);
  const [infrastructure, setInfrastructure] = useState<InfrastructureStepData>(defaultInfrastructure);
  const [painPoints, setPainPoints] = useState<PainPointsStepData>(defaultPainPoints);
  const [goals, setGoals] = useState<GoalsStepData>(defaultGoals);
  const [followUpResponses, setFollowUpResponses] = useState<Record<string, string>>({});
  const [isGeneratingAI, setIsGeneratingAI] = useState(false);

  const totalQuestionsInStep = STEP_QUESTION_COUNTS[currentStep] || 1;
  const isLastQuestionInStep = currentQuestion >= totalQuestionsInStep;
  const isFirstQuestionInStep = currentQuestion === 1;

  // API hooks
  const { data: options } = useEdTechOptions();
  const { data: profileData, refetch: refetchProfile } = useEdTechProfile(schoolId, open);
  const saveProfile = useSaveEdTechProfile();
  const generateAI = useGenerateAI();
  const submitFollowUp = useSubmitFollowUp();
  const completeProfile = useCompleteProfile();

  const profile = profileData?.profile;

  // Load existing profile data
  useEffect(() => {
    if (!profile) return;

    // eslint-disable-next-line react-hooks/set-state-in-effect
    setInfrastructure({
      totalDevices: profile.totalDevices,
      deviceTypes: profile.deviceTypes || defaultDeviceTypes,
      networkQuality: profile.networkQuality,
      internetSpeed: profile.internetSpeed,
      lmsPlatform: profile.lmsPlatform,
      existingSoftware: profile.existingSoftware || [],
      itStaffCount: profile.itStaffCount,
      deviceAge: profile.deviceAge,
    });
    setPainPoints({
      painPoints: profile.painPoints || [],
      supportSatisfaction: profile.supportSatisfaction,
      biggestChallenges: profile.biggestChallenges || [],
      supportFrequency: profile.supportFrequency,
      avgResolutionTime: profile.avgResolutionTime,
      biggestFrustration: profile.biggestFrustration,
      wishList: profile.wishList,
    });
    setGoals({
      strategicGoals: profile.strategicGoals || [],
      budgetRange: profile.budgetRange,
      timeline: profile.timeline,
      expansionPlans: profile.expansionPlans,
      priorityRanking: profile.priorityRanking || [],
      decisionMakers: profile.decisionMakers || [],
    });
    setFollowUpResponses(profile.followUpResponses || {});
  }, [profile]);

  // Reset state when modal closes
  useEffect(() => {
    if (!open) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setCurrentStep(1);
      setCurrentQuestion(1);
      setIsGeneratingAI(false);
    }
  }, [open]);

  // Reset question when step changes
  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setCurrentQuestion(1);
  }, [currentStep]);

  const handleNextQuestion = () => {
    if (isLastQuestionInStep) {
      return; // Should use Save & Continue instead
    }
    setCurrentQuestion((q) => q + 1);
  };

  const handlePrevQuestion = () => {
    if (isFirstQuestionInStep) {
      if (currentStep > 1) {
        setCurrentStep((s) => s - 1);
        // Set to last question of previous step
        setCurrentQuestion(STEP_QUESTION_COUNTS[currentStep - 1] || 1);
      }
      return;
    }
    setCurrentQuestion((q) => q - 1);
  };

  const handleSaveAndNext = useCallback(async () => {
    // Save current data
    const data = {
      schoolId,
      ...infrastructure,
      ...painPoints,
      ...goals,
    };

    try {
      await saveProfile.mutateAsync(data);

      if (currentStep < 4) {
        setCurrentStep(currentStep + 1);

        // When moving to step 4, generate AI if not already done
        if (currentStep === 3 && !profile?.aiSummary) {
          setIsGeneratingAI(true);
          // Need to get the profile ID after save
          const result = await refetchProfile();
          const savedProfile = result.data?.profile;
          if (savedProfile?.id) {
            await generateAI.mutateAsync({ profileId: savedProfile.id, schoolId });
            await refetchProfile();
          }
          setIsGeneratingAI(false);
        }
      }
    } catch (error) {
      console.error('Failed to save profile:', error);
      setIsGeneratingAI(false);
    }
  }, [currentStep, infrastructure, painPoints, goals, schoolId, profile, saveProfile, generateAI, refetchProfile]);

  const handleComplete = useCallback(async () => {
    if (!profile?.id) return;

    try {
      // Submit follow-up responses if any
      if (Object.keys(followUpResponses).length > 0) {
        await submitFollowUp.mutateAsync({
          profileId: profile.id,
          schoolId,
          data: { responses: followUpResponses },
        });
      }

      // Mark as complete
      await completeProfile.mutateAsync({ profileId: profile.id, schoolId });

      onClose();
    } catch (error) {
      console.error('Failed to complete profile:', error);
    }
  }, [profile, followUpResponses, schoolId, submitFollowUp, completeProfile, onClose]);

  const handleGenerateAI = useCallback(async () => {
    if (!profile?.id) return;

    setIsGeneratingAI(true);
    try {
      await generateAI.mutateAsync({ profileId: profile.id, schoolId });
      await refetchProfile();
    } catch (error) {
      console.error('Failed to generate AI:', error);
    }
    setIsGeneratingAI(false);
  }, [profile, schoolId, generateAI, refetchProfile]);

  const isSaving = saveProfile.isPending || generateAI.isPending || submitFollowUp.isPending || completeProfile.isPending;

  const renderStepContent = () => {
    switch (currentStep) {
      case 1:
        return (
          <InfrastructureStep
            data={infrastructure}
            options={options}
            onChange={setInfrastructure}
            currentQuestion={currentQuestion}
            totalQuestions={totalQuestionsInStep}
          />
        );
      case 2:
        return (
          <PainPointsStep
            data={painPoints}
            options={options}
            onChange={setPainPoints}
            currentQuestion={currentQuestion}
            totalQuestions={totalQuestionsInStep}
          />
        );
      case 3:
        return (
          <GoalsStep
            data={goals}
            options={options}
            onChange={setGoals}
            currentQuestion={currentQuestion}
            totalQuestions={totalQuestionsInStep}
          />
        );
      case 4:
        return (
          <AIFollowUpStep
            profile={profile || {} as EdTechProfile}
            followUpResponses={followUpResponses}
            onResponseChange={(id, response) =>
              setFollowUpResponses((prev) => ({ ...prev, [id]: response }))
            }
            isGenerating={isGeneratingAI}
          />
        );
      default:
        return null;
    }
  };

  return (
    <Modal open={open} onClose={onClose} className="max-w-3xl max-h-[90vh]">
      <ModalHeader onClose={onClose}>
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 text-white">
            <Sparkles className="h-4 w-4" />
          </div>
          <div>
            <h2 className="text-lg font-semibold">EdTech Profile Assessment</h2>
            <p className="text-sm text-gray-500">Help us understand your school's technology landscape</p>
          </div>
        </div>
      </ModalHeader>

      <ModalBody className="max-h-[60vh] overflow-y-auto">
        <AssessmentStepper
          currentStep={currentStep}
          onStepClick={setCurrentStep}
          allowNavigation={!!profile}
        />
        {renderStepContent()}
      </ModalBody>

      <ModalFooter>
        <div className="flex items-center justify-between w-full">
          <Button
            variant="outline"
            onClick={handlePrevQuestion}
            disabled={(currentStep === 1 && isFirstQuestionInStep) || isSaving}
          >
            <ChevronLeft className="w-4 h-4 mr-1" />
            Back
          </Button>

          <div className="flex gap-2">
            {currentStep === 4 && !profile?.aiSummary && !isGeneratingAI && (
              <Button
                variant="outline"
                onClick={handleGenerateAI}
                disabled={isSaving}
              >
                <Sparkles className="w-4 h-4 mr-1" />
                Generate AI Analysis
              </Button>
            )}

            {currentStep < 4 ? (
              isLastQuestionInStep ? (
                <Button
                  onClick={handleSaveAndNext}
                  disabled={isSaving}
                  className="bg-indigo-600 hover:bg-indigo-700"
                >
                  {isSaving ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-1 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    <>
                      Save & Continue
                      <ChevronRight className="w-4 h-4 ml-1" />
                    </>
                  )}
                </Button>
              ) : (
                <Button
                  onClick={handleNextQuestion}
                  className="bg-indigo-600 hover:bg-indigo-700"
                >
                  Next
                  <ChevronRight className="w-4 h-4 ml-1" />
                </Button>
              )
            ) : (
              <Button
                onClick={handleComplete}
                disabled={isSaving || isGeneratingAI}
                className="bg-green-600 hover:bg-green-700"
              >
                {isSaving ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-1 animate-spin" />
                    Completing...
                  </>
                ) : (
                  <>
                    <Check className="w-4 h-4 mr-1" />
                    Complete Assessment
                  </>
                )}
              </Button>
            )}
          </div>
        </div>
      </ModalFooter>
    </Modal>
  );
}
