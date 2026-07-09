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

func CreateCard(uid int, req dto.CreateCardRequest) (dto.CardResponse, error) {
	content, err := normalizeCardContent(req.Content)
	if err != nil {
		return dto.CardResponse{}, err
	}

	tagIDs := normalizeIDList(req.TagIds)
	for _, tagID := range tagIDs {
		if err := ensureTagAccessible(uid, tagID); err != nil {
			return dto.CardResponse{}, err
		}
	}

	card, err := database.CreateCard(model.Card{
		UserId:        uid,
		Title:         strings.TrimSpace(req.Title),
		Content:       content,
		Tags:          model.StringArray(normalizeStringList(req.Tags)),
		ExternalLinks: model.StringArray(normalizeStringList(req.ExternalLinks)),
		InternalLinks: model.StringArray(normalizeStringList(req.InternalLinks)),
	})
	if err != nil {
		return dto.CardResponse{}, err
	}

	if len(tagIDs) > 0 {
		if _, err := ReplaceCardTags(uid, card.Id, tagIDs); err != nil {
			return dto.CardResponse{}, err
		}
	}

	syncCardVectorBestEffort(card)
	return buildCardResponse(uid, card), nil
}

func GetCardByID(uid int, id int) (dto.CardResponse, error) {
	card, err := database.GetCardByIDAndUserID(id, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.CardResponse{}, vo.ErrCardNotFound
		}
		return dto.CardResponse{}, err
	}

	return buildCardResponse(uid, card), nil
}

func QueryCards(uid int, filter dto.QueryCardDTO) ([]dto.CardResponse, error) {
	page, pageSize := normalizePagination(filter.Page, filter.PageSize)

	cards, err := database.QueryCardsByUserID(database.CardQueryFilter{
		UserID:   uid,
		Content:  strings.TrimSpace(filter.Content),
		Tags:     normalizeStringList(filter.Tags),
		TagIDs:   normalizeIDList(filter.TagIds),
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	return toCardResponses(cards), nil
}

func UpdateCard(uid int, id int, req dto.UpdateCardRequest) (dto.CardResponse, error) {
	var tagIDs []int
	if req.TagIds != nil {
		tagIDs = normalizeIDList(*req.TagIds)
		for _, tagID := range tagIDs {
			if err := ensureTagAccessible(uid, tagID); err != nil {
				return dto.CardResponse{}, err
			}
		}
	}

	updates, err := buildCardUpdates(req)
	if err != nil {
		if !errors.Is(err, vo.ErrNoFieldsToUpdate) || req.TagIds == nil {
			return dto.CardResponse{}, err
		}
	}

	if len(updates) > 0 {
		if err := database.UpdateCardByIDAndUserID(id, uid, updates); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.CardResponse{}, vo.ErrCardNotFound
			}
			return dto.CardResponse{}, err
		}
	} else if err := ensureCardAccessible(uid, id); err != nil {
		return dto.CardResponse{}, err
	}

	if req.TagIds != nil {
		if _, err := ReplaceCardTags(uid, id, tagIDs); err != nil {
			return dto.CardResponse{}, err
		}
	}

	resp, err := GetCardByID(uid, id)
	if err != nil {
		return dto.CardResponse{}, err
	}

	if cardVectorFieldsChanged(updates) {
		syncCardVectorBestEffort(model.Card{
			Id:      resp.Id,
			UserId:  resp.UserId,
			Title:   resp.Title,
			Content: resp.Content,
		})
	}

	return resp, nil
}

func DeleteCard(uid int, id int) error {
	if err := database.SoftDeleteCardByIDAndUserID(id, uid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return vo.ErrCardNotFound
		}
		return err
	}
	if err := database.SoftDeleteCardTagsByCardIDAndUserID(uid, id); err != nil {
		return err
	}
	deleteCardVectorBestEffort(uid, id)
	return nil
}

func normalizeCardContent(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", vo.ErrCardContentRequired
	}
	return value, nil
}

func normalizeStringList(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{}, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
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

func buildCardUpdates(req dto.UpdateCardRequest) (map[string]any, error) {
	updates := make(map[string]any)

	if req.Title != nil {
		updates["title"] = strings.TrimSpace(*req.Title)
	}

	if req.Content != nil {
		content, err := normalizeCardContent(*req.Content)
		if err != nil {
			return nil, err
		}
		updates["content"] = content
	}

	if req.Tags != nil {
		updates["tags"] = model.StringArray(normalizeStringList(*req.Tags))
	}

	if req.ExternalLinks != nil {
		updates["external_links"] = model.StringArray(normalizeStringList(*req.ExternalLinks))
	}

	if req.InternalLinks != nil {
		updates["internal_links"] = model.StringArray(normalizeStringList(*req.InternalLinks))
	}

	if len(updates) == 0 {
		return nil, vo.ErrNoFieldsToUpdate
	}

	return updates, nil
}

func toCardResponses(cards []model.Card) []dto.CardResponse {
	resp := make([]dto.CardResponse, 0, len(cards))
	for _, card := range cards {
		resp = append(resp, buildCardResponse(card.UserId, card))
	}
	return resp
}

func buildCardResponse(uid int, card model.Card) dto.CardResponse {
	resp := toCardResponse(card)

	cardTags, err := database.QueryCardTagsByUserID(database.CardTagQueryFilter{
		UserID: uid,
		CardID: card.Id,
	})
	if err != nil {
		return resp
	}

	resp.TagIds = make([]int, 0, len(cardTags))
	for _, cardTag := range cardTags {
		resp.TagIds = append(resp.TagIds, cardTag.TagId)
	}

	return resp
}

func toCardResponse(card model.Card) dto.CardResponse {
	return dto.CardResponse{
		Id:            card.Id,
		UserId:        card.UserId,
		Title:         card.Title,
		Content:       card.Content,
		Tags:          []string(card.Tags),
		ExternalLinks: []string(card.ExternalLinks),
		InternalLinks: []string(card.InternalLinks),
		CreatedAt:     card.CreatedAt,
		UpdatedAt:     card.UpdatedAt,
	}
}
