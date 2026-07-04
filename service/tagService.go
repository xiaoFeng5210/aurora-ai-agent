package service

import (
	"aurora-agent/database"
	"aurora-agent/database/model"
	"aurora-agent/handler/dto"
	"aurora-agent/handler/vo"
	"errors"
	"strings"

	"gorm.io/gorm"
)

func CreateTag(uid int, req dto.CreateTagRequest) (dto.TagResponse, error) {
	name, err := normalizeTagName(req.Name)
	if err != nil {
		return dto.TagResponse{}, err
	}

	exists, err := database.TagNameExists(uid, name, 0)
	if err != nil {
		return dto.TagResponse{}, err
	}
	if exists {
		return dto.TagResponse{}, vo.ErrTagNameExists
	}

	tag, err := database.CreateTag(model.Tag{
		UserId: uid,
		Name:   name,
	})
	if err != nil {
		return dto.TagResponse{}, err
	}

	return toTagResponse(tag), nil
}

func GetTagByID(uid int, id int) (dto.TagResponse, error) {
	tag, err := database.GetTagByIDAndUserID(id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TagResponse{}, vo.ErrTagNotFound
		}
		return dto.TagResponse{}, err
	}

	return toTagResponse(tag), nil
}

func QueryTags(uid int, filter dto.QueryTagDTO) ([]dto.TagResponse, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)

	tags, err := database.QueryTagsByUserID(database.TagQueryFilter{
		UserID:   uid,
		Name:     strings.TrimSpace(filter.Name),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	return toTagResponses(tags), nil
}

func UpdateTag(uid int, id int, req dto.UpdateTagRequest) (dto.TagResponse, error) {
	updates, err := buildTagUpdates(uid, id, req)
	if err != nil {
		return dto.TagResponse{}, err
	}

	if err := database.UpdateTagByIDAndUserID(id, uid, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TagResponse{}, vo.ErrTagNotFound
		}
		return dto.TagResponse{}, err
	}

	return GetTagByID(uid, id)
}

func DeleteTag(uid int, id int) error {
	if err := database.SoftDeleteTagByIDAndUserID(id, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrTagNotFound
		}
		return err
	}

	return database.SoftDeleteCardTagsByTagIDAndUserID(uid, id)
}

func AddTagToCard(uid int, cardID int, tagID int) (dto.CardTagResponse, error) {
	if err := ensureCardAccessible(uid, cardID); err != nil {
		return dto.CardTagResponse{}, err
	}
	if err := ensureTagAccessible(uid, tagID); err != nil {
		return dto.CardTagResponse{}, err
	}

	cardTag, err := database.CreateCardTag(model.CardTag{
		UserId: uid,
		CardId: cardID,
		TagId:  tagID,
	})
	if err != nil {
		return dto.CardTagResponse{}, err
	}

	return toCardTagResponse(cardTag), nil
}

func ReplaceCardTags(uid int, cardID int, tagIDs []int) ([]dto.CardTagResponse, error) {
	if err := ensureCardAccessible(uid, cardID); err != nil {
		return nil, err
	}

	tagIDs = normalizeIDList(tagIDs)
	for _, tagID := range tagIDs {
		if err := ensureTagAccessible(uid, tagID); err != nil {
			return nil, err
		}
	}

	cardTags, err := database.ReplaceCardTagsByCardIDAndUserID(uid, cardID, tagIDs)
	if err != nil {
		return nil, err
	}

	return toCardTagResponses(cardTags), nil
}

func DeleteTagFromCard(uid int, cardID int, tagID int) error {
	if err := ensureCardAccessible(uid, cardID); err != nil {
		return err
	}
	if err := ensureTagAccessible(uid, tagID); err != nil {
		return err
	}

	if err := database.SoftDeleteCardTagByCardIDAndTagID(uid, cardID, tagID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrCardTagNotFound
		}
		return err
	}
	return nil
}

func QueryCardTags(uid int, filter dto.QueryCardTagDTO) ([]dto.CardTagResponse, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)

	if filter.CardId > 0 {
		if err := ensureCardAccessible(uid, filter.CardId); err != nil {
			return nil, err
		}
	}
	if filter.TagId > 0 {
		if err := ensureTagAccessible(uid, filter.TagId); err != nil {
			return nil, err
		}
	}

	cardTags, err := database.QueryCardTagsByUserID(database.CardTagQueryFilter{
		UserID:   uid,
		CardID:   filter.CardId,
		TagID:    filter.TagId,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	return toCardTagResponses(cardTags), nil
}

func QueryTagsByCard(uid int, cardID int) ([]dto.TagResponse, error) {
	if err := ensureCardAccessible(uid, cardID); err != nil {
		return nil, err
	}

	cardTags, err := database.QueryCardTagsByUserID(database.CardTagQueryFilter{
		UserID: uid,
		CardID: cardID,
	})
	if err != nil {
		return nil, err
	}

	tags := make([]dto.TagResponse, 0, len(cardTags))
	for _, cardTag := range cardTags {
		tag, err := database.GetTagByIDAndUserID(cardTag.TagId, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, err
		}
		tags = append(tags, toTagResponse(tag))
	}

	return tags, nil
}

func normalizeTagName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", vo.ErrTagNameRequired
	}
	return value, nil
}

func buildTagUpdates(uid int, id int, req dto.UpdateTagRequest) (map[string]any, error) {
	updates := make(map[string]any)

	if req.Name != nil {
		name, err := normalizeTagName(*req.Name)
		if err != nil {
			return nil, err
		}

		exists, err := database.TagNameExists(uid, name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, vo.ErrTagNameExists
		}

		updates["name"] = name
	}

	if len(updates) == 0 {
		return nil, vo.ErrNoFieldsToUpdate
	}

	return updates, nil
}

func ensureCardAccessible(uid int, cardID int) error {
	if _, err := database.GetCardByIDAndUserID(cardID, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrCardNotFound
		}
		return err
	}
	return nil
}

func ensureTagAccessible(uid int, tagID int) error {
	if _, err := database.GetTagByIDAndUserID(tagID, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrTagNotFound
		}
		return err
	}
	return nil
}

func normalizeIDList(values []int) []int {
	if len(values) == 0 {
		return []int{}
	}

	seen := make(map[int]struct{}, len(values))
	normalized := make([]int, 0, len(values))
	for _, value := range values {
		if value <= 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}

	return normalized
}

func toTagResponses(tags []model.Tag) []dto.TagResponse {
	resp := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		resp = append(resp, toTagResponse(tag))
	}
	return resp
}

func toTagResponse(tag model.Tag) dto.TagResponse {
	return dto.TagResponse{
		Id:        tag.Id,
		UserId:    tag.UserId,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}
}

func toCardTagResponses(cardTags []model.CardTag) []dto.CardTagResponse {
	resp := make([]dto.CardTagResponse, 0, len(cardTags))
	for _, cardTag := range cardTags {
		resp = append(resp, toCardTagResponse(cardTag))
	}
	return resp
}

func toCardTagResponse(cardTag model.CardTag) dto.CardTagResponse {
	return dto.CardTagResponse{
		Id:        cardTag.Id,
		UserId:    cardTag.UserId,
		CardId:    cardTag.CardId,
		TagId:     cardTag.TagId,
		CreatedAt: cardTag.CreatedAt,
		UpdatedAt: cardTag.UpdatedAt,
	}
}
