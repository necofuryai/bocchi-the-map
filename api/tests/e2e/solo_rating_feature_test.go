package e2e_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/masyusakai/bocchi-the-map/api/internal/domain/entities"
	"github.com/masyusakai/bocchi-the-map/api/pb"
	"github.com/masyusakai/bocchi-the-map/api/tests/helpers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSoloRatingFeature(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Solo Rating Feature Suite")
}

// BDD E2E Test - Outside Loop
var _ = Describe("Solo Rating Feature", func() {
	var (
		suite    *helpers.CommonTestSuite
		authUser *entities.User
		testSpot *entities.Spot
	)

	BeforeEach(func() {
		suite = helpers.NewCommonTestSuite()
		authUser = suite.CreateTestUser()
		testSpot = suite.CreateTestSpot()
	})

	AfterEach(func() {
		suite.Cleanup()
	})

	Context("Given I am an authenticated solo traveler", func() {
		BeforeEach(func() {
			suite.AuthenticateUser(authUser)
		})

		Context("When I rate a spot for solo-friendliness", func() {
			It("Then the rating should be saved and reflected in spot statistics", func() {
				By("Creating a solo-friendly rating request")
				ratingRequest := &pb.CreateRatingRequest{
					SpotId:             testSpot.ID,
					UserId:             authUser.ID,
					SoloFriendlyRating: 5,
					Categories: []string{
						"quiet_atmosphere",
						"wifi_available",
						"single_seating",
					},
					Comment: "Perfect spot for solo work with great WiFi",
				}

				By("Submitting the rating via gRPC")
				response, err := suite.SpotClient.CreateRating(context.Background(), ratingRequest)

				By("Verifying the response")
				Expect(err).ToNot(HaveOccurred())
				Expect(response.Rating).ToNot(BeNil())
				Expect(response.Rating.Id).ToNot(BeEmpty())
				Expect(response.Rating.SoloFriendlyRating).To(Equal(int32(5)))
				Expect(response.Rating.UserId).To(Equal(authUser.ID))

				By("Verifying the rating appears in spot statistics")
				spotResponse, err := suite.SpotClient.GetSpot(context.Background(), &pb.GetSpotRequest{
					Id: testSpot.ID,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(spotResponse.Spot.SoloFriendlyStats.AverageRating).To(BeNumerically(">=", 4.5))
				Expect(spotResponse.Spot.SoloFriendlyStats.TotalRatings).To(Equal(int32(1)))
			})
		})

		Context("When I rate a spot that I've already rated", func() {
			BeforeEach(func() {
				// Create initial rating
				initialRating := &pb.CreateRatingRequest{
					SpotId:             testSpot.ID,
					UserId:             authUser.ID,
					SoloFriendlyRating: 3,
					Comment:            "Initial rating",
				}
				_, err := suite.SpotClient.CreateRating(context.Background(), initialRating)
				Expect(err).ToNot(HaveOccurred())
			})

			It("Then the rating should be updated, not duplicated", func() {
				By("Creating an updated rating request")
				updatedRating := &pb.CreateRatingRequest{
					SpotId:             testSpot.ID,
					UserId:             authUser.ID,
					SoloFriendlyRating: 5,
					Comment:            "Updated rating after revisit",
				}

				By("Submitting the updated rating")
				response, err := suite.SpotClient.CreateRating(context.Background(), updatedRating)

				By("Verifying the response")
				Expect(err).ToNot(HaveOccurred())
				Expect(response.Rating.SoloFriendlyRating).To(Equal(int32(5)))
				Expect(response.Rating.Comment).To(Equal("Updated rating after revisit"))

				By("Verifying only one rating exists for this user-spot combination")
				ratingsResponse, err := suite.SpotClient.GetSpotRatings(context.Background(), &pb.GetSpotRatingsRequest{
					SpotId: testSpot.ID,
				})
				Expect(err).ToNot(HaveOccurred())
				
				userRatings := 0
				for _, rating := range ratingsResponse.Ratings {
					if rating.UserId == authUser.ID {
						userRatings++
					}
				}
				Expect(userRatings).To(Equal(1))
			})
		})

		Context("When I submit an invalid rating", func() {
			It("Then I should receive a validation error", func() {
				By("Creating an invalid rating request")
				invalidRating := &pb.CreateRatingRequest{
					SpotId:             "", // Empty spot ID
					UserId:             authUser.ID,
					SoloFriendlyRating: 10, // Invalid rating (should be 1-5)
					Comment:            "",
				}

				By("Submitting the invalid rating")
				_, err := suite.SpotClient.CreateRating(context.Background(), invalidRating)

				By("Verifying the validation error")
				Expect(err).To(HaveOccurred())
				grpcStatus, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(grpcStatus.Code()).To(Equal(codes.InvalidArgument))
				Expect(grpcStatus.Message()).To(ContainSubstring("validation failed"))
			})
		})
	})

	Context("Given I am not authenticated", func() {
		It("Then I should not be able to create ratings", func() {
			By("Creating a rating request without authentication")
			ratingRequest := &pb.CreateRatingRequest{
				SpotId:             testSpot.ID,
				UserId:             "unauthenticated-user",
				SoloFriendlyRating: 5,
				Comment:            "Attempting to rate without auth",
			}

			By("Submitting the rating")
			_, err := suite.SpotClient.CreateRating(context.Background(), ratingRequest)

			By("Verifying the authentication error")
			Expect(err).To(HaveOccurred())
			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.Unauthenticated))
		})
	})
})